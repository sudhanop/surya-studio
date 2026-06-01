/* =============================================================================
   Surya Photography — SQL Server Database Schema
   Run via create-database.bat or open in SSMS and execute (F5)
   ============================================================================= */

SET NOCOUNT ON;
SET QUOTED_IDENTIFIER ON;
SET ANSI_NULLS ON;
GO

-- Create database (skip if already exists)
IF NOT EXISTS (SELECT 1 FROM sys.databases WHERE name = N'SuryaPhotography')
BEGIN
    CREATE DATABASE SuryaPhotography;
END
GO

USE SuryaPhotography;
GO

/* ---------------------------------------------------------------------------
   Drop existing objects (development reset — comment out in production)
   --------------------------------------------------------------------------- */
IF OBJECT_ID(N'dbo.function_event_dates', N'U') IS NOT NULL DROP TABLE dbo.function_event_dates;
IF OBJECT_ID(N'dbo.functions', N'U') IS NOT NULL DROP TABLE dbo.functions;
IF OBJECT_ID(N'dbo.site_settings', N'U') IS NOT NULL DROP TABLE dbo.site_settings;
IF OBJECT_ID(N'dbo.portfolio_media', N'U') IS NOT NULL DROP TABLE dbo.portfolio_media;
IF OBJECT_ID(N'dbo.inquiries', N'U') IS NOT NULL DROP TABLE dbo.inquiries;
IF OBJECT_ID(N'dbo.categories', N'U') IS NOT NULL DROP TABLE dbo.categories;
IF OBJECT_ID(N'dbo.admin', N'U') IS NOT NULL DROP TABLE dbo.admin;
GO

/* ---------------------------------------------------------------------------
   admin — single admin account (JWT + bcrypt in application)
   --------------------------------------------------------------------------- */
CREATE TABLE dbo.admin (
    id              BIGINT IDENTITY(1,1) NOT NULL,
    username        NVARCHAR(100) NOT NULL,
    password_hash   NVARCHAR(255) NOT NULL,
    email           NVARCHAR(255) NOT NULL,
    display_name    NVARCHAR(150) NULL,
    is_active       BIT NOT NULL CONSTRAINT DF_admin_is_active DEFAULT (1),
    last_login_at   DATETIME2(3) NULL,
    created_at      DATETIME2(3) NOT NULL CONSTRAINT DF_admin_created_at DEFAULT (SYSUTCDATETIME()),
    updated_at      DATETIME2(3) NOT NULL CONSTRAINT DF_admin_updated_at DEFAULT (SYSUTCDATETIME()),
    CONSTRAINT PK_admin PRIMARY KEY CLUSTERED (id),
    CONSTRAINT UQ_admin_username UNIQUE (username),
    CONSTRAINT UQ_admin_email UNIQUE (email)
);
GO

/* ---------------------------------------------------------------------------
   categories — portfolio categories
   --------------------------------------------------------------------------- */
CREATE TABLE dbo.categories (
    id              BIGINT IDENTITY(1,1) NOT NULL,
    name            NVARCHAR(100) NOT NULL,
    slug            NVARCHAR(100) NOT NULL,
    description     NVARCHAR(500) NULL,
    cover_image     NVARCHAR(500) NULL,
    display_order   INT NOT NULL CONSTRAINT DF_categories_display_order DEFAULT (0),
    is_active       BIT NOT NULL CONSTRAINT DF_categories_is_active DEFAULT (1),
    created_at      DATETIME2(3) NOT NULL CONSTRAINT DF_categories_created_at DEFAULT (SYSUTCDATETIME()),
    updated_at      DATETIME2(3) NOT NULL CONSTRAINT DF_categories_updated_at DEFAULT (SYSUTCDATETIME()),
    CONSTRAINT PK_categories PRIMARY KEY CLUSTERED (id),
    CONSTRAINT UQ_categories_slug UNIQUE (slug)
);
GO

CREATE NONCLUSTERED INDEX IX_categories_active_order
    ON dbo.categories (is_active, display_order);
GO

/* ---------------------------------------------------------------------------
   portfolio_media — photos and videos per category
   --------------------------------------------------------------------------- */
CREATE TABLE dbo.portfolio_media (
    id              BIGINT IDENTITY(1,1) NOT NULL,
    category_id     BIGINT NOT NULL,
    title           NVARCHAR(200) NULL,
    caption         NVARCHAR(500) NULL,
    media_type      NVARCHAR(20) NOT NULL,  -- 'photo' | 'video'
    file_path       NVARCHAR(500) NOT NULL,
    thumbnail_path  NVARCHAR(500) NULL,
    mime_type       NVARCHAR(100) NULL,
    file_size_bytes BIGINT NULL,
    duration_sec    INT NULL,
    is_featured     BIT NOT NULL CONSTRAINT DF_portfolio_media_featured DEFAULT (0),
    display_order   INT NOT NULL CONSTRAINT DF_portfolio_media_order DEFAULT (0),
    is_published    BIT NOT NULL CONSTRAINT DF_portfolio_media_published DEFAULT (1),
    created_at      DATETIME2(3) NOT NULL CONSTRAINT DF_portfolio_media_created_at DEFAULT (SYSUTCDATETIME()),
    updated_at      DATETIME2(3) NOT NULL CONSTRAINT DF_portfolio_media_updated_at DEFAULT (SYSUTCDATETIME()),
    CONSTRAINT PK_portfolio_media PRIMARY KEY CLUSTERED (id),
    CONSTRAINT FK_portfolio_media_category
        FOREIGN KEY (category_id) REFERENCES dbo.categories (id) ON DELETE CASCADE,
    CONSTRAINT CK_portfolio_media_type CHECK (media_type IN (N'photo', N'video'))
);
GO

CREATE NONCLUSTERED INDEX IX_portfolio_media_category
    ON dbo.portfolio_media (category_id, media_type, is_published, display_order);
GO

CREATE NONCLUSTERED INDEX IX_portfolio_media_featured
    ON dbo.portfolio_media (is_featured, is_published)
    WHERE is_featured = 1;
GO

/* ---------------------------------------------------------------------------
   inquiries — public contact form submissions
   --------------------------------------------------------------------------- */
CREATE TABLE dbo.inquiries (
    id              BIGINT IDENTITY(1,1) NOT NULL,
    customer_name   NVARCHAR(150) NOT NULL,
    phone_number    NVARCHAR(30) NOT NULL,
    occasion_type   NVARCHAR(100) NOT NULL,
    wanted_date     DATE NULL,
    address         NVARCHAR(500) NULL,
    message         NVARCHAR(MAX) NULL,
    status          NVARCHAR(30) NOT NULL CONSTRAINT DF_inquiries_status DEFAULT (N'new'),
    contacted_at    DATETIME2(3) NULL,
    created_at      DATETIME2(3) NOT NULL CONSTRAINT DF_inquiries_created_at DEFAULT (SYSUTCDATETIME()),
    updated_at      DATETIME2(3) NOT NULL CONSTRAINT DF_inquiries_updated_at DEFAULT (SYSUTCDATETIME()),
    CONSTRAINT PK_inquiries PRIMARY KEY CLUSTERED (id),
    CONSTRAINT CK_inquiries_status CHECK (status IN (N'new', N'contacted', N'converted', N'archived'))
);
GO

CREATE NONCLUSTERED INDEX IX_inquiries_status_created
    ON dbo.inquiries (status, created_at DESC);
GO

/* ---------------------------------------------------------------------------
   functions — internal studio production tracker
   --------------------------------------------------------------------------- */
CREATE TABLE dbo.functions (
    id                  BIGINT IDENTITY(1,1) NOT NULL,
    inquiry_id          BIGINT NULL,
    customer_name       NVARCHAR(150) NOT NULL,
    phone_number        NVARCHAR(30) NOT NULL,
    address             NVARCHAR(500) NULL,
    function_type       NVARCHAR(100) NOT NULL,
    function_date       DATE NOT NULL,
    total_amount        DECIMAL(18,2) NOT NULL CONSTRAINT DF_functions_total DEFAULT (0),
    advance_paid        DECIMAL(18,2) NOT NULL CONSTRAINT DF_functions_advance DEFAULT (0),
    balance_amount      AS (total_amount - advance_paid) PERSISTED,
    assigned_editor     NVARCHAR(150) NULL,
    assigned_date       DATE NULL,
    album_status        NVARCHAR(50) NOT NULL CONSTRAINT DF_functions_album_status DEFAULT (N'not_started'),
    video_status        NVARCHAR(50) NOT NULL CONSTRAINT DF_functions_video_status DEFAULT (N'not_started'),
    customer_booking_notes NVARCHAR(MAX) NULL,
    services_json       NVARCHAR(MAX) NULL,
    complimentary_json  NVARCHAR(MAX) NULL,
    delivery_status     NVARCHAR(50) NOT NULL CONSTRAINT DF_functions_delivery_status DEFAULT (N'pending'),
    overall_status      NVARCHAR(30) NOT NULL CONSTRAINT DF_functions_overall_status DEFAULT (N'upcoming'),
    admin_notes         NVARCHAR(MAX) NULL,
    drive_links         NVARCHAR(MAX) NULL,
    created_at          DATETIME2(3) NOT NULL CONSTRAINT DF_functions_created_at DEFAULT (SYSUTCDATETIME()),
    updated_at          DATETIME2(3) NOT NULL CONSTRAINT DF_functions_updated_at DEFAULT (SYSUTCDATETIME()),
    CONSTRAINT PK_functions PRIMARY KEY CLUSTERED (id),
    CONSTRAINT FK_functions_inquiry
        FOREIGN KEY (inquiry_id) REFERENCES dbo.inquiries (id) ON DELETE SET NULL,
    CONSTRAINT CK_functions_overall_status CHECK (
        overall_status IN (N'upcoming', N'completed', N'editing', N'album_ready', N'delivered')
    )
);
GO

CREATE NONCLUSTERED INDEX IX_functions_function_date
    ON dbo.functions (function_date, overall_status);
GO

CREATE NONCLUSTERED INDEX IX_functions_overall_status
    ON dbo.functions (overall_status, function_date);
GO

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
GO

CREATE NONCLUSTERED INDEX IX_function_event_dates_function
    ON dbo.function_event_dates (function_id, sort_order, event_date);
GO

CREATE TABLE dbo.site_settings (
    setting_key     NVARCHAR(100) NOT NULL,
    setting_value   NVARCHAR(MAX) NULL,
    updated_at      DATETIME2(3) NOT NULL CONSTRAINT DF_site_settings_updated DEFAULT (SYSUTCDATETIME()),
    CONSTRAINT PK_site_settings PRIMARY KEY CLUSTERED (setting_key)
);
GO

/* ---------------------------------------------------------------------------
   Seed: portfolio categories
   --------------------------------------------------------------------------- */
INSERT INTO dbo.categories (name, slug, description, display_order) VALUES
    (N'Wedding',           N'wedding',           N'Wedding photography and cinematic films', 1),
    (N'Baby Shower',       N'baby-shower',       N'Baby shower celebrations', 2),
    (N'Puberty',           N'puberty',           N'Puberty ceremony coverage', 3),
    (N'Birthday',          N'birthday',          N'Birthday events and parties', 4),
    (N'Ear Piercing',      N'ear-piercing',      N'Ear piercing ceremonies', 5),
    (N'Couple Shoots',     N'couple-shoots',     N'Couple and pre-wedding shoots', 6),
    (N'Maternity',         N'maternity',         N'Maternity photography sessions', 7),
    (N'Outdoor',           N'outdoor',           N'Outdoor and location shoots', 8),
    (N'Temple Functions',  N'temple-functions',  N'Temple and religious functions', 9);
GO

/* ---------------------------------------------------------------------------
   Seed: default admin
   Username: surya@admin.com
   Password: surya@1995  (bcrypt hash below)
   --------------------------------------------------------------------------- */
INSERT INTO dbo.admin (username, password_hash, email, display_name)
VALUES (
    N'surya@admin.com',
    N'$2a$10$rNXUCcV9I.3mAswLIutSu.Yae/04qfLrZWjvPYGuMdSOF6u4qPVru',
    N'surya@admin.com',
    N'Surya Photography Admin'
);
GO

INSERT INTO dbo.site_settings (setting_key, setting_value) VALUES
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
    (N'whatsapp', N'919715241568');
GO

PRINT 'SuryaPhotography database schema created successfully.';
PRINT 'Default admin — username: surya@admin.com | password: surya@1995';
GO
