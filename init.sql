begin;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE extension IF NOT EXISTS btree_gist;

create table client (
    uid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    surname text not null,
    name text not null,
    patronymic text not null,
    email text not null,
    login text not null,
    password_hash text not null,
    role_uid uuid not null,
    role text not null,
    approved bool not null default false
);

create table access_group (
	uid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	name text not null,
	description text not null
);


create table equipment (
	uid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	name text not null,
	description text not null,
	type text not null,
	manufacturer text,
	model text,
	room text not null
);

create table equipment_schedule(
	uid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	equipment_uid uuid not null,
	time_interval tsrange not null,
	maintaince_flag bool not null default false
);


create table inventory (
	uid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	name text not null,
	description text not null,
	type text not null,
	manufacturer text,
	quantity decimal not null,
	unit text not null
);

create table experiment (
	uid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	name text not null,
	description text not null,
	start_ts timestamp not null,
	end_ts timestamp not null
);

create table inventory_in_access_group(
	access_group_uid uuid not null,
	inventory_uid uuid not null
);

create table equipment_in_access_group(
	access_group_uid uuid not null,
	equipment_uid uuid not null
);


create table clients_in_access_group(
	access_group_uid uuid not null,
	client_uid uuid not null
);

create table clients_in_experiment(
	experiment_uid uuid not null,
	client_uid uuid not null
);

create table inventory_in_experiment(
	experiment_uid uuid not null,
	inventory_uid uuid not null,
	quantity decimal not null
);

create table equipment_schedule_in_experiment(
	experiment_uid uuid not null,
	equipment_schedule_uid uuid not null
);

create table maintaince(
    uid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text not null,
    description text not null,
    start_ts timestamp not null,
    end_ts timestamp not null
);

create table equipment_schedule_in_maintaince(
    maintaince_uid uuid not null,
    equipment_schedule_uid  uuid not null
);

create table clients_in_maintaince(
  maintaince_uid uuid not null,
  client_uid  uuid not null
);



alter table clients_in_access_group
	add constraint fk_clients_in_access_group_access_group foreign key(access_group_uid) references access_group(uid);

alter table clients_in_access_group
	add constraint fk_clients_in_access_group_client foreign key(client_uid) references client(uid);



alter table equipment_in_access_group
	add constraint fk_equipment_in_access_group_access_group foreign key(access_group_uid) references access_group(uid);

alter table equipment_in_access_group
	add constraint fk_equipment_in_access_group_equipment foreign key(equipment_uid) references equipment(uid);



alter table inventory_in_access_group
	add constraint fk_inventory_in_access_group_access_group foreign key(access_group_uid) references access_group(uid);

alter table inventory_in_access_group
	add constraint fk_inventory_in_access_group_inventory foreign key(inventory_uid) references inventory(uid);



create unique index on client(login);


alter table equipment_schedule
	add constraint fk_equipment_schedule_equipment foreign key(equipment_uid) references equipment(uid);

alter table equipment_schedule
	add constraint unique_equipment_schedule_interval exclude using gist  (maintaince_flag with =, equipment_uid WITH =, time_interval WITH &&);



alter table clients_in_experiment
	add constraint fk_clients_in_experiment_experiment foreign key(experiment_uid) references experiment(uid);

alter table clients_in_experiment
	add constraint fk_clients_in_experiment_client foreign key(client_uid) references client(uid);



alter table inventory_in_experiment
	add constraint fk_inventory_in_experiment_experiment foreign key(experiment_uid) references experiment(uid);

alter table inventory_in_experiment
	add constraint fk_inventory_in_experiment_inventory foreign key(inventory_uid) references inventory(uid);



alter table equipment_schedule_in_experiment
	add constraint fk_equipment_schedule_in_experiment_experiment foreign key(experiment_uid) references experiment(uid);

alter table equipment_schedule_in_experiment
	add constraint fk_equipment_schedule_in_experiment_equipment_schedule foreign key(equipment_schedule_uid) references equipment_schedule(uid);



create unique index on equipment(name);


create unique index on inventory(name);
alter table inventory
    add constraint positive_quantity check(quantity >= 0);


create unique index on access_group(name);



alter table equipment_schedule_in_maintaince
    add constraint fk_equipment_schedule_in_maintaince_maintaince foreign key(maintaince_uid) references maintaince(uid);

alter table equipment_schedule_in_maintaince
    add constraint fk_equipment_schedule_in_maintaince_equipment_schedule foreign key(equipment_schedule_uid) references equipment_schedule(uid);



alter table clients_in_maintaince
    add constraint fk_clients_in_maintaince_maintaince foreign key(maintaince_uid) references maintaince(uid);

alter table clients_in_maintaince
    add constraint fk_clients_in_maintaince_clients foreign key(client_uid) references client(uid);

commit;

alter table maintaince
    add column end_ts timestamp not null;

