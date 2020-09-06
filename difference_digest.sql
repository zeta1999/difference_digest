CREATE OR REPLACE FUNCTION pg_temp.f_hash(idx int, element bigint) RETURNS numeric
    RETURNS NULL ON NULL INPUT
    IMMUTABLE
    LANGUAGE plpythonu
AS $$
  hashSeeds = [18269989962351869307, 9143901319630896501, 2072764263930962169, 417226483919003799, 16485935163296413021]
  mask = int("0x5555555555555555", 16)

  inner = (mask * (element ^ (element >> 32))) % 18446744073709551616
  return (hashSeeds[idx] * inner) % 18446744073709551616
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

CREATE OR REPLACE AGGREGATE pg_temp.f_bit_xor (bigint) (
    SFUNC = pg_temp.f_bit_xor_sfunc,
    STYPE = bigint,
    INITCOND = 0
);

CREATE OR REPLACE FUNCTION pg_temp.f_bit_xor_numeric_sfunc(agg numeric, value numeric) RETURNS numeric
  RETURNS NULL ON NULL INPUT
  IMMUTABLE
  LANGUAGE plpythonu
AS $$
  return int(agg) ^ int(value)
$$;

CREATE OR REPLACE AGGREGATE pg_temp.f_bit_xor_numeric (numeric) (
    SFUNC = pg_temp.f_bit_xor_numeric_sfunc,
    STYPE = numeric,
    INITCOND = 0
);
