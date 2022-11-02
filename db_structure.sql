create table files
(
    file_name text,
    blob      bytea,
    filesize  bigint,
    sent      timestamp(0),
    sentby    bigint
);

alter table files
    owner to nyrdyxoc;

create table messages
(
    id             bigserial,
    text           text,
    sent           timestamp(0),
    sentby         bigint,
    tel_chat_id    bigint,
    tel_message_id bigint,
    message_type   numeric,
    viewedby       integer,
    viewedat       timestamp with time zone,
    replyto        numeric
);

alter table messages
    owner to nyrdyxoc;

create index idx_messages_telchatid
    on messages (tel_chat_id);

create index idx_messages_viewedat
    on messages (viewedat);

create table buttons
(
    id               bigserial,
    rid              bigint,
    icon             varchar,
    button_values_id bigint,
    version          bigint,
    isactive         boolean,
    inserted         information_schema.time_stamp,
    insertedby       varchar,
    updated          timestamp,
    updatedby        varchar
);

alter table buttons
    owner to nyrdyxoc;

create table button_values
(
    id         bigserial,
    text       varchar,
    lang       varchar not null,
    inserted   information_schema.time_stamp,
    insertedby varchar,
    updated    information_schema.time_stamp,
    updatedby  varchar
);

alter table button_values
    owner to nyrdyxoc;

create table requests
(
    reqnumber          numeric,
    reqfrom            numeric,
    reqtype            varchar(200),
    datetime           timestamp with time zone default CURRENT_TIMESTAMP not null,
    id                 bigserial,
    status             integer                  default 0                 not null,
    servicesrequestsid bigint,
    feedback           varchar(50)
);

comment on column requests.status is 'sorğunun bütün sualları cavablandırlıb-1 , əks halda-0';

alter table requests
    owner to nyrdyxoc;

create table questions
(
    id                       bigserial,
    state                    numeric,
    question_type_id         numeric,
    request_text             varchar,
    request_error_text       varchar,
    response_validation_type varchar,
    response_type            integer
);

comment on column questions.response_type is '1-no list 2-list with only one selection 3-list with multi selection';

alter table questions
    owner to nyrdyxoc;

create table question_type
(
    id   bigserial,
    name varchar not null
);

alter table question_type
    owner to nyrdyxoc;

create table logs
(
    id        bigserial,
    chat_id   numeric,
    timestamp timestamp with time zone default CURRENT_TIMESTAMP,
    text      text,
    type      varchar
);

alter table logs
    owner to nyrdyxoc;

create table question_answers
(
    id             bigserial,
    questions_id   bigint not null,
    value          varchar,
    chat_id        bigint,
    timestamp      timestamp with time zone default CURRENT_TIMESTAMP,
    request_number bigint
);

alter table question_answers
    owner to nyrdyxoc;

create table request_statuses
(
    id         bigserial,
    status     varchar                                            not null,
    request_id bigint                                             not null,
    insertedby varchar                                            not null,
    insertedat timestamp with time zone default CURRENT_TIMESTAMP not null
);

alter table request_statuses
    owner to nyrdyxoc;

create table question_list
(
    id          bigserial,
    "order"     bigint,
    question_id bigint  not null,
    value       varchar not null
);

alter table question_list
    owner to nyrdyxoc;

create table voices
(
    voice       bytea,
    chatid      varchar(50),
    messageid   varchar(50),
    voicesize   numeric,
    duration    numeric,
    sentdate    timestamp with time zone,
    id          serial,
    messages_id numeric
);

alter table voices
    owner to nyrdyxoc;

create table servicesrequests
(
    id              serial,
    service_name    varchar(1000),
    request_type_id bigint
);

alter table servicesrequests
    owner to nyrdyxoc;

create table servicerequestscomponents
(
    services_requests_id  numeric,
    order_num             numeric,
    component_description varchar(100),
    component_type        varchar(50),
    data_driven           numeric(1),
    id                    serial
);

alter table servicerequestscomponents
    owner to nyrdyxoc;

create table servicerequestscomponentsdetails
(
    servicerequestscomponents_id numeric,
    component_id                 varchar(1000),
    component_name               varchar(1000),
    component_value              varchar(1000),
    component_label              varchar(1000),
    component_requiredsize       varchar(1000),
    component_placeholder        varchar(1000),
    component_minlength          varchar(1000),
    component_maxlength          varchar(1000),
    component_title              varchar(1000),
    component_mindate            varchar(1000),
    component_maxdate            varchar(1000),
    id                           serial
);

alter table servicerequestscomponentsdetails
    owner to nyrdyxoc;

create table usersservicesrequests
(
    tel_chat_id         varchar(100),
    servicesrequests_id numeric,
    request_datetime    timestamp with time zone
);

alter table usersservicesrequests
    owner to nyrdyxoc;

create table servicerequestscomponentsdatas
(
    servicerequestscomponents_id numeric,
    data_value                   varchar(1000),
    requests_id                  bigint
);

alter table servicerequestscomponentsdatas
    owner to nyrdyxoc;

create table test
(
    id   numeric,
    name varchar(40)
);

alter table test
    owner to nyrdyxoc;

create table request_type
(
    id   bigserial,
    name varchar(1000)
);

alter table request_type
    owner to nyrdyxoc;

