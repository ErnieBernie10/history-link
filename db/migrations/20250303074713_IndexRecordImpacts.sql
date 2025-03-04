-- migrate:up
create index idx_record_impacts on impact (record_id);

create index idx_record_history_record_id on record_history (record_id);

create index idx_impact_history_impact_id on impact_history (impact_id);

-- migrate:down
drop index idx_record_impacts;

drop index idx_record_history_record_id;

drop index idx_impact_history_impact_id;
