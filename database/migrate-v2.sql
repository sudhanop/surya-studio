/* Surya Photography — v2 migration (run on existing SuryaPhotography DB) */
USE SuryaPhotography;
GO

IF OBJECT_ID(N'dbo.function_event_dates', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.function_event_dates (
        id              BIGINT IDENTITY(1,1) NOT NULL,
        function_id     BIGINT NOT NULL,
        event_date      DATE NOT NULL,
        day_label       NVARCHAR(100) NULL,
        sort_order      INT NOT NULL CONSTRAINT DF_function_event_dates_order DEFAULT (0),
        CONSTRAINT PK_function_event_dates PRIMARY KEY CLUSTERED (id),
        CONSTRAINT FK_function_event_dates_function
            FOREIGN KEY (function_id) REFERENCES dbo.functions (id) ON DELETE CASCADE
    );
    CREATE NONCLUSTERED INDEX IX_function_event_dates_function
        ON dbo.function_event_dates (function_id, sort_order, event_date);
END
GO

IF COL_LENGTH('dbo.functions', 'customer_booking_notes') IS NULL
    ALTER TABLE dbo.functions ADD customer_booking_notes NVARCHAR(MAX) NULL;
GO
IF COL_LENGTH('dbo.functions', 'services_json') IS NULL
    ALTER TABLE dbo.functions ADD services_json NVARCHAR(MAX) NULL;
GO
IF COL_LENGTH('dbo.functions', 'complimentary_json') IS NULL
    ALTER TABLE dbo.functions ADD complimentary_json NVARCHAR(MAX) NULL;
GO

UPDATE dbo.functions SET album_status = N'not_started' WHERE album_status = N'pending';
UPDATE dbo.functions SET video_status = N'not_started' WHERE video_status = N'pending';
GO

IF OBJECT_ID(N'dbo.site_settings', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.site_settings (
        setting_key     NVARCHAR(100) NOT NULL,
        setting_value   NVARCHAR(MAX) NULL,
        updated_at      DATETIME2(3) NOT NULL CONSTRAINT DF_site_settings_updated DEFAULT (SYSUTCDATETIME()),
        CONSTRAINT PK_site_settings PRIMARY KEY CLUSTERED (setting_key)
    );
END
GO

MERGE dbo.site_settings AS t
USING (VALUES
    (N'events_covered', N'500+'),
    (N'years_of_craft', N'10+'),
    (N'happy_families', N'1000+'),
    (N'contact_email', N'suryaphotographyrsp@gmail.com'),
    (N'phone_primary', N'9715241568'),
    (N'phone_secondary', N'8884897499'),
    (N'instagram_url', N'https://www.instagram.com/surya_photography_nkl'),
    (N'youtube_url', N'https://www.youtube.com/@suryaphotography4303'),
    (N'address', N'Surya Photography, near DNC (Chamundi) Theater, opposite to Adhisindha Thirumana Mandabam, Pattanam road, Rasipuram'),
    (N'pincode', N'637408'),
    (N'whatsapp', N'919715241568')
) AS s(setting_key, setting_value)
ON t.setting_key = s.setting_key
WHEN NOT MATCHED THEN
    INSERT (setting_key, setting_value) VALUES (s.setting_key, s.setting_value);
GO

PRINT 'Migration v2 completed.';
GO
