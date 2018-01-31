--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.10
-- Dumped by pg_dump version 9.5.10

-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE submittal_log_facts (
    id bigint NOT NULL,
    company_id bigint,
    cost_code_id bigint,
    created_by_id bigint,
    location_id bigint,
    project_id bigint,
    received_from_id bigint,
    responsible_contractor_id bigint,
    specification_section_id bigint,
    submittal_log_id bigint,
    submittal_log_status_id bigint,
    submittal_manager_id bigint,
    submittal_package_id bigint,
    cost_code text,
    created_by text,
    custom_textarea_1 text,
    custom_textfield_1 text,
    deleted_at timestamp without time zone,
    date_received date,
    date_required_on_site date,
    date_issue date,
    date_created_at date,
    date_submit_by date,
    date_distributed date,
    date_final_due date,
    date_actual_delivery date,
    date_confirmed_delivery date,
    date_anticipated_delivery date,
    description text,
    design_team_review_time integer,
    internal_review_time integer,
    lead_time integer,
    location text,
    number text,
    package_name text,
    package_number text,
    planned_internal_review_completed_date date,
    planned_return_date date,
    planned_submit_by_date date,
    project_address text,
    project_bid_type text,
    project_city text,
    project_county text,
    project_date_created date,
    project_department jsonb,
    project_estimated_start_date date,
    project_estimated_completion_date date,
    project_description text,
    project_designated_market_area text,
    project_name text,
    project_notes text,
    project_number text,
    project_office text,
    project_owner_type text,
    project_parent_job text,
    project_phone text,
    project_program text,
    project_square_feet integer,
    project_region text,
    project_stage text,
    project_state text,
    project_type text,
    project_zip text,
    private boolean,
    received_from text,
    responsible_contractor text,
    revision text,
    scheduled_task text,
    spec_section_description text,
    spec_section_number text,
    status text,
    status_name text,
    submittal_manager text,
    submittal_type text,
    title text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    approver_names text[] DEFAULT '{}'::text[],
    approver_ids bigint[] DEFAULT '{}'::bigint[],
    approver_vendor_ids bigint[] DEFAULT '{}'::bigint[],
    approver_responses text[] DEFAULT '{}'::text[],
    approver_sent_dates date[] DEFAULT '{}'::date[],
    approver_returned_dates date[] DEFAULT '{}'::date[],
    approver_due_dates date[] DEFAULT '{}'::date[],
    responsed_vendor_ids bigint[] DEFAULT '{}'::bigint[],
    ball_in_court_names text[] DEFAULT '{}'::text[],
    ball_in_court_ids bigint[] DEFAULT '{}'::bigint[],
    ball_in_court_due_date date[] DEFAULT '{}'::date[]
);


ALTER TABLE ONLY submittal_log_facts
    ADD CONSTRAINT submittal_log_facts_pkey PRIMARY KEY (id);


--
-- Name: index_submittal_log_facts_on_project_id; Type: INDEX; Schema: public; Owner: williamhuang
--

CREATE INDEX index_submittal_log_facts_on_project_id ON submittal_log_facts USING btree (project_id);


--
-- Name: index_submittal_log_facts_on_submittal_log_id; Type: INDEX; Schema: public; Owner: williamhuang
--

CREATE UNIQUE INDEX index_submittal_log_facts_on_submittal_log_id ON submittal_log_facts USING btree (submittal_log_id);


--
-- PostgreSQL database dump complete
--
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE IF EXISTS submittal_log_facts CASCADE;
