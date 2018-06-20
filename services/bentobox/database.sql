--
-- PostgreSQL database dump
--

-- Dumped from database version 10.4 (Debian 10.4-2.pgdg90+1)
-- Dumped by pg_dump version 10.4 (Debian 10.4-2.pgdg90+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: blocknumberoftx; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.blocknumberoftx (
    inserted_ts bigint,
    block_id text,
    number_of_txs bigint
);


ALTER TABLE public.blocknumberoftx OWNER TO postgres;

--
-- Name: blocktx; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.blocktx (
    inserted_ts bigint,
    block_id text,
    tx_id text
);


ALTER TABLE public.blocktx OWNER TO postgres;

--
-- Name: ethdata; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ethdata (
    inserted_ts bigint,
    kind text,
    hash text,
    cid text,
    value text,
    last_ipfs_add_ts bigint,
    ipfs_success_ts bigint
);


ALTER TABLE public.ethdata OWNER TO postgres;

--
-- Name: lastblock; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.lastblock (
    inserted_ts bigint,
    number_id bigint
);


ALTER TABLE public.lastblock OWNER TO postgres;

--
-- Name: txreceipts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.txreceipts (
    inserted_ts bigint,
    tx_id text,
    tx_receipts_id text
);


ALTER TABLE public.txreceipts OWNER TO postgres;

--
-- Name: wantfromdevp2p; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.wantfromdevp2p (
    inserted_ts bigint,
    kind text,
    key text,
    last_request_ts bigint,
    success_ts bigint
);


ALTER TABLE public.wantfromdevp2p OWNER TO postgres;

--
-- Name: block_id_bnot_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX block_id_bnot_idx ON public.blocknumberoftx USING btree (block_id);


--
-- Name: block_id_bt_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX block_id_bt_idx ON public.blocktx USING btree (block_id);


--
-- Name: cid_ed_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX cid_ed_idx ON public.ethdata USING btree (cid);


--
-- Name: hash_ed_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX hash_ed_idx ON public.ethdata USING btree (hash);


--
-- Name: inserted_ts_bnot_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX inserted_ts_bnot_idx ON public.blocknumberoftx USING btree (inserted_ts);


--
-- Name: inserted_ts_bt_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX inserted_ts_bt_idx ON public.blocktx USING btree (inserted_ts);


--
-- Name: inserted_ts_ed_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX inserted_ts_ed_idx ON public.ethdata USING btree (inserted_ts);


--
-- Name: inserted_ts_lb_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX inserted_ts_lb_idx ON public.lastblock USING btree (inserted_ts);


--
-- Name: inserted_ts_tr_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX inserted_ts_tr_idx ON public.txreceipts USING btree (inserted_ts);


--
-- Name: inserted_ts_wfd_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX inserted_ts_wfd_idx ON public.wantfromdevp2p USING btree (inserted_ts);


--
-- Name: ipfs_success_ts_ed_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ipfs_success_ts_ed_idx ON public.ethdata USING btree (ipfs_success_ts);


--
-- Name: key_wfd_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX key_wfd_idx ON public.wantfromdevp2p USING btree (key);


--
-- Name: last_ipfs_add_ts_ed_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX last_ipfs_add_ts_ed_idx ON public.ethdata USING btree (last_ipfs_add_ts);


--
-- Name: last_request_ts_wfd_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX last_request_ts_wfd_idx ON public.wantfromdevp2p USING btree (last_request_ts);


--
-- Name: number_id_lb_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX number_id_lb_idx ON public.lastblock USING btree (number_id);


--
-- Name: success_ts_wfd_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX success_ts_wfd_idx ON public.wantfromdevp2p USING btree (success_ts);


--
-- Name: tx_id_bt_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX tx_id_bt_idx ON public.blocktx USING btree (tx_id);


--
-- Name: tx_id_tr_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX tx_id_tr_idx ON public.txreceipts USING btree (tx_id);


--
-- Name: tx_receipts_id_tr_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX tx_receipts_id_tr_idx ON public.txreceipts USING btree (tx_receipts_id);


--
-- PostgreSQL database dump complete
--

