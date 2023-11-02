DROP TABLE IF EXISTS passwords;

DROP TABLE IF EXISTS notes;

DROP TABLE IF EXISTS cards; 

SELECT lo_unlink(binaries.bin_id) FROM binaries;
DROP TABLE IF EXISTS binaries;

DROP TABLE IF EXISTS users;