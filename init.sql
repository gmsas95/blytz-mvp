-- Create databases for each microservice
CREATE DATABASE auth_db;
CREATE DATABASE products_db;
CREATE DATABASE auction_db;
CREATE DATABASE chat_db;
CREATE DATABASE orders_db;
CREATE DATABASE payments_db;
CREATE DATABASE logistics_db;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE auth_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE products_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE auction_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE chat_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE orders_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE payments_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE logistics_db TO postgres;