ALTER TABLE sch1.ShoppingCarts ADD CONSTRAINT FKShoppingCartsCustomers FOREIGN KEY(CustomerId, CustomerName)
    REFERENCES sch1.Customers(CustomerId, CustomerName) ON DELETE CASCADE