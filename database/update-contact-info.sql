/* Force studio contact details everywhere in the database */
USE SuryaPhotography;
GO

MERGE dbo.site_settings AS t
USING (VALUES
    (N'contact_email', N'suryaphotographyrsp@gmail.com'),
    (N'phone_primary', N'9715241568'),
    (N'phone_secondary', N'8884897499'),
    (N'whatsapp', N'919715241568'),
    (N'instagram_url', N'https://www.instagram.com/surya_photography_nkl'),
    (N'youtube_url', N'https://www.youtube.com/@suryaphotography4303'),
    (N'address', N'Surya Photography, near DNC (Chamundi) Theater, opposite to Adhisindha Thirumana Mandabam, Pattanam road, Rasipuram'),
    (N'pincode', N'637408')
) AS s(setting_key, setting_value)
ON t.setting_key = s.setting_key
WHEN MATCHED THEN
    UPDATE SET setting_value = s.setting_value, updated_at = SYSUTCDATETIME()
WHEN NOT MATCHED THEN
    INSERT (setting_key, setting_value) VALUES (s.setting_key, s.setting_value);
GO

PRINT 'Contact info updated in site_settings.';
GO
