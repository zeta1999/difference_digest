CREATE EXTENSION IF NOT EXISTS plpythonu;

CREATE OR REPLACE FUNCTION pg_temp.f_hash(idx int, element bigint) RETURNS bigint
    RETURNS NULL ON NULL INPUT
    IMMUTABLE
    LANGUAGE plpythonu
AS $$
  hashSeeds = [1826998997, 914390139, 207279169, 4179003799, 1648963021]
  mask = int("0x555555555555", 16)

  inner = (mask * (element ^ (element >> 32))) % 4294967296
  return (hashSeeds[idx] * inner) % 4294967296
$$; 

CREATE OR REPLACE FUNCTION pg_temp.f_trailing_zeros(element numeric) RETURNS int
    RETURNS NULL ON NULL INPUT
    IMMUTABLE
    LANGUAGE plpythonu
AS $$
  x = int(element)
  count = 0
  if x == 0:
    return count

  while ((x & 1) == 0): 
      x = x >> 1
      count += 1
      
  return count 
$$; 

CREATE OR REPLACE FUNCTION pg_temp.f_bit_xor_sfunc(agg bigint, value bigint) RETURNS bigint
  RETURNS NULL ON NULL INPUT
  IMMUTABLE
  LANGUAGE SQL
AS $$
SELECT agg # value
$$;

CREATE AGGREGATE pg_temp.f_bit_xor (bigint) (
    SFUNC = pg_temp.f_bit_xor_sfunc,
    STYPE = bigint,
    INITCOND = 0
);
