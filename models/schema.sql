CREATE TABLE PackageMetadata (
    ID VARCHAR(120),
    NAME VARCHAR(100),
    Version VARCHAR(20),
    License VARCHAR(20),
    ReadMe VARCHAR(1000),
    RepoURL TINYTEXT,
    NetScore DOUBLE,
    BusFactor DOUBLE,
    Correctness DOUBLE,
    RampUp DOUBLE,
    ResponsiveMaintainer DOUBLE,
    LicenseScore DOUBLE,
    GoodPinningPractice DOUBLE,
    GoodEngineeringProcess DOUBLE
);