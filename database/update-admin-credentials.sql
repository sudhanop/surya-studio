/* Update admin login to surya@admin.com / surya@1995 (run against existing SuryaPhotography DB) */
USE SuryaPhotography;
GO

UPDATE dbo.admin
SET
    username = N'surya@admin.com',
    password_hash = N'$2a$10$rNXUCcV9I.3mAswLIutSu.Yae/04qfLrZWjvPYGuMdSOF6u4qPVru',
    email = N'surya@admin.com',
    display_name = N'Surya Photography Admin',
    updated_at = SYSUTCDATETIME()
WHERE username IN (N'admin', N'surya@admin.com')
   OR id = (SELECT MIN(id) FROM dbo.admin);
GO

PRINT 'Admin credentials updated — username: surya@admin.com | password: surya@1995';
GO
