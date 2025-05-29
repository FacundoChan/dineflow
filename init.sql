CREATE DATABASE IF NOT EXISTS dineflow;

USE dineflow;

DROP TABLE IF EXISTS `order_stock`;
DROP TABLE IF EXISTS `product_images`;

CREATE TABLE `order_stock` (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  product_id VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  quantity INT UNSIGNED NOT NULL DEFAULT 0,
  price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
  description TEXT,
  version INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `product_images` (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  product_id VARCHAR(255) NOT NULL,
  img_url VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (product_id) REFERENCES order_stock(product_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO order_stock(product_id, name, quantity, price, description)
VALUES
  ('prod_S3CrGrzAS1MZsK','Cake', 20, 10, 'Description of cake, 蛋糕的介绍'),
  ('prod_S3Cr3l2WHdiL53', 'Beef', 50, 20, 'Description of Beef, 牛排的介绍');

INSERT INTO product_images(product_id, img_url)
VALUES 
  ('prod_S3CrGrzAS1MZsK', 'https://raw.gitmirror.com/FChanAux/web-pic-bed/main/image/202504110105963.webp'),
  ('prod_S3CrGrzAS1MZsK', 'https://raw.gitmirror.com/FChanAux/web-pic-bed/main/image/202504110105964.webp'),
  ('prod_S3CrGrzAS1MZsK', 'https://raw.gitmirror.com/FChanAux/web-pic-bed/main/image/202504110105961.webp'),
  ('prod_S3Cr3l2WHdiL53', 'https://raw.gitmirror.com/FChanAux/web-pic-bed/main/image/202504110105966.webp'),
  ('prod_S3Cr3l2WHdiL53', 'https://raw.gitmirror.com/FChanAux/web-pic-bed/main/image/202504110105965.webp');
