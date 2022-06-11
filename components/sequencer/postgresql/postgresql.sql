//
// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.



    -- Table: public.layotto_incr

-- DROP TABLE IF EXISTS public.layotto_incr;

CREATE TABLE IF NOT EXISTS public.layotto_incr
(
    id bigint NOT NULL,
    value_id bigint NOT NULL,
    biz_tag character(255) COLLATE pg_catalog."default" NOT NULL,
    create_time bigint,
    update_time bigint,
    CONSTRAINT layotto_incr_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.layotto_incr
    OWNER to postgres;

INSERT INTO public.layotto_incr(
    id, value_id, biz_tag, create_time, update_time)
VALUES (?, ?, ?, ?, ?);
