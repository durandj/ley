INSERT INTO Users (
    ID,
    Username,
    Status,
    CreatedOn
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING ID, Username, Status, CreatedOn, ModifiedOn
;
