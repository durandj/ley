SELECT
    ID,
    Username,
    Status,
    CreatedOn,
    ModifiedOn
FROM Users
WHERE
    username = $1
LIMIT 1
;
