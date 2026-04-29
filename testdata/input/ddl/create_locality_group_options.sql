CREATE LOCALITY GROUP spill_to_hdd
OPTIONS (storage = 'ssd', ssd_to_hdd_spill_timespan = '10d')