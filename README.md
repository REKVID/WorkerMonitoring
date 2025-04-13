
# TaskBot - Telegram Bot for Managing Photo Reports

## 📝 Description
TgTaskBot is a Telegram bot designed to automate the process of collecting and managing photo reports. 
The bot allows workers to send photos from job sites, and administrators to manage and view received reports.

## 🚀 Key Features

### 👷 For Workers:
- Selection of work site
- Specifying job position
- Sending photo reports
- Automatic saving of photos with metadata

### 👨‍💼 For Administrators:
- Viewing all uploaded photos
- Filtering photos by date/site/position
- Managing the list of sites and positions
- Viewing the file system
- Access to system logs

## 🛠 Technical Requirements

- Windows/Linux/MacOS
- Go 1.16 or higher
- Access to Telegram Bot API
- Permissions to create files and folders

## 📂 Project Structure

```bash
TgTaskBot/
├── main.go                 # Application entry point
├── configuration.json      # Configuration file
├── Config/                 # Configuration package
│   └── config.go
├── internal/               # Internal logic
│   ├── admin/              # Admin functionality
│   ├── bot/                # Core bot logic
│   └── worker/             # Worker functionality
├── pkg/                    # Shared packages
│   ├── handlers/           # Handlers
│   └── logger/             # Logging
```

## ⚙️ Installation and Setup

1. Clone the repository:
```bash
git clone https://github.com/your-username/TgTaskBot.git
cd TgTaskBot
```

2. Create the configuration file `Config/configuration.json`:
```json
{
    "bot_token": "YOUR_BOT_TOKEN",
    "admin_key": "YOUR_ADMIN_KEY",
    "positions": ["Security Guard", "Cleaner"],
    "objects": ["Site 1", "Site 2"],
    "base_dir": "./photos",
    "log_path": "./app.log"
}
```

3. Build the project:
```bash
# For Windows
go build -o TgTaskBot.exe main.go

# For Linux/MacOS
GOOS=windows GOARCH=amd64 go build -o TgTaskBot.exe main.go
```

## 🚀 Running the Bot

1. Via batch file:
   - Run `start.bat`

2. Directly:
   - Run the executable file `TgTaskBot.exe`

## 📱 Using the Bot

### Getting Started:
1. Send the `/start` command
2. Select a role (Worker/Admin)

### For Workers:
1. Select the "Worker" role
2. Choose a site from the list
3. Specify your position
4. Send a photo

### For Administrators:
1. Select the "Admin" role
2. Enter the admin key
3. Use available options:
   - View photos
   - View file system
   - View logs
   - Add sites/positions

## 📁 Photo Storage Structure

```
photos/
├── [Position]/
│   ├── [Site]/
│   │   └── [Date]/
│   │       └── [Username].jpg
```

## 🔒 Security
- Admin panel access is password-protected
- All actions are logged
- Photos are stored in an organized structure

## 📝 Logging
- All actions are recorded in `app.log`
- Logs include:
  - User actions
  - System errors
  - Photo uploads
  - Administrative actions

## ⚠️ Error Handling
- Automatic creation of required directories
- Permission checks
- Input data validation
- Network error handling

## 🔧 Troubleshooting

### Common Issues:

1. **Configuration Errors**
   - Check for the presence of `configuration.json`
   - Ensure the bot token is correct

2. **Access Errors**
   - Check file/folder creation permissions
   - Ensure all necessary directories exist

3. **Photo Upload Problems**
   - Check internet connection
   - Ensure photo size does not exceed Telegram limits

## Support
If issues arise:
1. Check the log file `app.log`
2. Verify the configuration
3. Check file system access permissions

## 🔄 Updates
- Regularly check the repository for updates
- Backup `configuration.json` before updating

## 📜 License
MIT License - free to use and modify

## Author
https://github.com/REKVID
