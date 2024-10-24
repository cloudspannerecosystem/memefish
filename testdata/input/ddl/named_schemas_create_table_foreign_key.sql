CREATE TABLE sch1.ShoppingCarts (
  CartId INT64 NOT NULL,
  CustomerId INT64 NOT NULL,
  CustomerName STRING(MAX) NOT NULL,
  CONSTRAINT FKShoppingCartsCustomers FOREIGN KEY(CustomerId, CustomerName)
    REFERENCES sch1.Customers(CustomerId, CustomerName) ON DELETE CASCADE,
) PRIMARY KEY(CartId)