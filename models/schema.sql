CREATE TABLE PkgMetadata (
    ID VARCHAR(222) NOT NULL,
    PRIMARY KEY(ID),
    NAME VARCHAR(214) NOT NULL,
    Version VARCHAR(8) NOT NULL,
    License VARCHAR(20),
    ReadMe VARCHAR(400),
    RepoURL TINYTEXT,
    NetScore DOUBLE NOT NULL,
    BusFactor DOUBLE NOT NULL,
    Correctness DOUBLE NOT NULL,
    RampUp DOUBLE NOT NULL,
    ResponsiveMaintainer DOUBLE NOT NULL,
    LicenseScore DOUBLE NOT NULL,
    GoodPinningPractice DOUBLE NOT NULL,
    GoodEngineeringProcess DOUBLE NOT NULL,
    Date TIMESTAMP,
    Action VARCHAR(10)
);