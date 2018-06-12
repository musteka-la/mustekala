const waterfall = require('async/waterfall')
const IPFS = require('ipfs')
const Repo = require('ipfs-repo')
const CID = require('cids')
const base32 = require('base32.js')
const mh = require('multihashes')
const ConcatStream = require('concat-stream')
const http = require('http')
const FsStore = require('datastore-fs')

const ETH_PROTOCOL = process.env.ETH_PROTOCOL || 'http'
const ETH_HOST = process.env.ETH_HOST || 'localhost'
const ETH_PORT = process.env.ETH_PORT || '5001'
const uriBase = `${ETH_PROTOCOL}://${ETH_HOST}:${ETH_PORT}/api/v0/block/get?arg=`

console.log(`Mounting Parity as data store: ${ETH_PROTOCOL}://${ETH_HOST}:${ETH_PORT}`)

class CustomFsStore extends FsStore {
  get (key, cb) {
    // console.log('CustomFsStore - get', arguments)
    const cid = interceptEthCid(key)
    if (!cid) return super.get.apply(this, arguments)
    fetchByCid(cid, cb)
  }
  has (key, cb) {
    // console.log('CustomFsStore - has', arguments)
    // console.log('CustomFsStore - has', key.toString())
    const cid = interceptEthCid(key)
    if (!cid) return super.has.apply(this, arguments)
    fetchByCid(cid, (err, result) => {
      if (err) return cb(null, false)
      cb(null, true)
    })
  }
}

function interceptEthCid(key){
  // console.log('interceptEthCid', key.toString())
  const slug = key.toString()
  if (['/SHARDING'].includes(slug)) return
  const rawKey = slug.split('/')[2]
  if (!rawKey) return
  const cid = dsKeyToCid(rawKey)
  // console.log(cid)
  const isEth = cid.codec.includes('eth-') || cid.codec === 'raw'
  if (!isEth) return
  return cid
}

function dsKeyToCid(rawKey) {
  // console.log('rawKey:', rawKey)
  const decoder = new base32.Decoder()
  const buf = decoder.write(rawKey).finalize()
  // console.log('buf:', buf)
  const cid = new CID(buf)
  return cid
}

waterfall([
  (cb) => prepareIpfs(cb),
  (node, cb) => setupHttpApi(node)
])

function prepareIpfs(cb) {

  const repo = new Repo('./ipfs-repo', {
    storageBackends: {
      blocks: CustomFsStore,
    },
  })
  const node = new IPFS({
    repo: repo,
    start: true,
    Bootstrap: []
  })

  node.once('ready', () => cb(null, node))
  node.on('error', (err) => {
    console.error(err)
  }) // Node has hit some error while initing/starting
  node.on('ready', () => console.log('ipfs node ready'))     // Node has successfully finished initing the repo
  node.on('init', () => console.log('ipfs node init'))     // Node has successfully finished initing the repo
  node.on('start', () => console.log('ipfs node start'))    // Node has started
  node.on('stop', () => console.log('ipfs node stop'))     // Node has stopped
}


function fetchByCid(cid, cb) {
  // filter for valid ethereum hashes
  if (mh.decode(cid.multihash).name !== 'keccak-256') return cb(new Error('Parity fetch failed - unsupported hash type'))
  // continue fetching
  // hot fix for https://github.com/paritytech/parity/issues/4172#issuecomment-314722992
  const fixedCid = (cid.codec !== 'eth-storage-trie') ? cid : new CID(cid.version, 'eth-state-trie', cid.multihash)
  const cidString = fixedCid.toBaseEncodedString()
  const uri = uriBase + cidString
  http.get(uri, (res) => {
    res.pipe(ConcatStream((result) => {
      if (res.statusCode !== 200) return cb(new Error(`Parity fetch failed - ${res.statusCode} - ${result}`))
      cb(null, result)
    }))
    res.once('error', cb)
  })
}

function noopRead(cid, cb) {
  cb()
}

function setupHttpApi(node) {
  // const HttpAPI = require('ipfs/src/http')
  console.log('Starting http server....')
  httpAPI = new HttpApiServer(node)
  httpAPI.start((err) => {
    console.log('http server is up!')
    if (err && err.code === 'ENOENT') {
      console.log('Error: no ipfs repo found in ' + repoPath)
      console.log('please run: jsipfs init')
      process.exit(1)
    }
    if (err) {
      throw err
    }
    console.log('Daemon is ready')
  })

  process.on('SIGINT', () => {
    console.log('Received interrupt signal, shutting down..')
    httpAPI.stop((err) => {
      if (err) {
        throw err
      }
      process.exit(0)
    })
  })
}

//
// Http API Server
//

const series = require('async/series')
const Hapi = require('hapi')
const errorHandler = require('ipfs/src/http/error-handler')
const multiaddr = require('multiaddr')
const setHeader = require('hapi-set-header')

function uriToMultiaddr (uri) {
  const ipPort = uri.split('/')[2].split(':')
  return `/ip4/${ipPort[0]}/tcp/${ipPort[1]}`
}

class HttpApiServer {

  constructor (node) {
    this.node = node
    this.log = console.log
    this.log.error = console.error
  }

  start (cb) {
    series([
      (cb) => {
        this.log('fetching config')
        this.node._repo.config.get((err, config) => {
          if (err) {
            return callback(err)
          }

          // CORS is enabled by default
          this.server = new Hapi.Server({
            connections: { routes: { cors: true } },
            debug: {
              log: '*',
              request: ['received'],
            },
          })

          this.server.app.ipfs = this.node
          const api = config.Addresses.API.split('/')
          const gateway = config.Addresses.Gateway.split('/')

          // select which connection with server.select(<label>) to add routes
          this.server.connection({
            // host: api[2],
            // host: '0.0.0.0',
            port: api[4],
            labels: 'API'
          })

          this.server.connection({
            // host: gateway[2],
            // host: '0.0.0.0',
            port: gateway[4],
            labels: 'Gateway'
          })

          // log connection details
          console.log({
            // host: api[2],
            // host: '0.0.0.0',
            port: api[4],
            labels: 'API'
          })
          console.log({
            // host: gateway[2],
            // host: '0.0.0.0',
            port: gateway[4],
            labels: 'Gateway'
          })

          // Nicer errors
          errorHandler(this, this.server)

          // load routes
          // require('./routes')(this.server)
          require('ipfs/src/http/api/routes')(this.server)
          require('ipfs/src/http/gateway/routes')(this.server)

          // Set default headers
          setHeader(this.server,
            'Access-Control-Allow-Headers',
            'X-Stream-Output, X-Chunked-Output, X-Content-Length')
          setHeader(this.server,
            'Access-Control-Expose-Headers',
            'X-Stream-Output, X-Chunked-Output, X-Content-Length')

          this.server.start(cb)
        })
      },
      (cb) => {
        const api = this.server.select('API')
        const gateway = this.server.select('Gateway')
        this.apiMultiaddr = multiaddr('/ip4/127.0.0.1/tcp/' + api.info.port)
        api.info.ma = uriToMultiaddr(api.info.uri)
        gateway.info.ma = uriToMultiaddr(gateway.info.uri)

        this.log('API is listening on: %s', api.info.ma)
        this.log('Gateway (readonly) is listening on: %s', gateway.info.ma)

        // for the CLI to know the where abouts of the API
        this.node._repo.apiAddr.set(api.info.ma, cb)
      }
    ], cb)

    this.stop = (callback) => {
      this.log('stopping')
      series([
        (cb) => this.server.stop(cb),
        (cb) => this.node.stop(cb)
      ], (err) => {
        if (err) {
          this.log.error(err)
          this.log('There were errors stopping')
        }
        callback()
      })
    }
  }

}