ALTER TABLE tasks
    ADD COLUMN target_file_name                varchar(30),
    ADD COLUMN target_file_validation          boolean default false,
    ADD COLUMN target_file_validation_language varchar(10);