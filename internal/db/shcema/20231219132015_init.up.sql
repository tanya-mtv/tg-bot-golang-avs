IF OBJECT_ID(N'TGUnit', N'U') IS NULL
CREATE TABLE TGUnit
(
   id INT PRIMARY KEY IDENTITY (1,1),
   ModifiedDate DATETIME,
   Owner INT NOT NULL REFERENCES Product (id),
   TGStatus INT,
   Barcode NVARCHAR(100),
   Volume DECIMAL(16,8),
   Weight DECIMAL(16,8),
   Unit INT NOT NULL REFERENCES UnitQualifier (id)
);

IF OBJECT_ID(N'TGUsers', N'U') IS NULL
CREATE TABLE TGUsers
(
   id INT PRIMARY KEY IDENTITY (1,1),
   ModifiedDate DATETIME,
   Name NVARCHAR(100),
   TGID NVARCHAR(100),
   IsAdmin BIT,
   IsActive BIT,
);
