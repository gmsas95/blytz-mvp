#!/usr/bin/env python3
"""
Simple HTTP server to serve the Blytz Auction MVP frontend
This allows testing the frontend without Docker dependencies
"""

import http.server
import socketserver
import os
import sys
from urllib.parse import urlparse
import json

class CustomHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        """Handle GET requests"""
        # Parse the URL
        parsed_path = urlparse(self.path)
        path = parsed_path.path

        # Serve the frontend
        if path == '/' or path == '/frontend' or path == '/frontend/':
            # Redirect to the actual frontend
            self.send_response(302)
            self.send_header('Location', '/frontend/index.html')
            self.end_headers()
            return

        if path.startswith('/frontend/'):
            # Serve frontend files
            file_path = path[1:]  # Remove leading slash
            if os.path.exists(file_path):
                self.serve_file(file_path)
            else:
                self.send_error(404, f"File not found: {file_path}")
            return

        # Default behavior for other files
        super().do_GET()

    def serve_file(self, file_path):
        """Serve a file with proper content type"""
        try:
            # Determine content type
            if file_path.endswith('.html'):
                content_type = 'text/html'
            elif file_path.endswith('.css'):
                content_type = 'text/css'
            elif file_path.endswith('.js'):
                content_type = 'application/javascript'
            else:
                content_type = 'application/octet-stream'

            # Read and serve the file
            with open(file_path, 'rb') as f:
                content = f.read()

            self.send_response(200)
            self.send_header('Content-Type', content_type)
            self.send_header('Access-Control-Allow-Origin', '*')
            self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
            self.send_header('Access-Control-Allow-Headers', 'Content-Type, Authorization')
            self.send_header('Content-Length', str(len(content)))
            self.end_headers()
            self.wfile.write(content)

        except FileNotFoundError:
            self.send_error(404, f"File not found: {file_path}")
        except Exception as e:
            self.send_error(500, f"Server error: {str(e)}")

    def do_OPTIONS(self):
        """Handle OPTIONS requests for CORS"""
        self.send_response(200)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type, Authorization')
        self.send_header('Access-Control-Max-Age', '86400')
        self.end_headers()

    def log_message(self, format, *args):
        """Custom logging"""
        print(f"[{self.date_time_string()}] {format % args}")

def main():
    """Main function"""
    # Change to the project directory
    project_dir = os.path.dirname(os.path.abspath(__file__))
    os.chdir(project_dir)

    # Configuration
    PORT = 8081  # Changed from 8080 to avoid conflict
    HOST = '0.0.0.0'

    print("🚀 Starting Blytz Auction MVP Simple Server")
    print("==========================================")
    print(f"Server will run on http://localhost:{PORT}")
    print(f"Frontend available at: http://localhost:{PORT}/frontend/index.html")
    print("")
    print("📋 Features:")
    print("- Serves the frontend HTML/CSS/JS files")
    print("- Handles CORS for API calls")
    print("- Static file serving")
    print("")
    print("💡 Usage:")
    print("1. Open http://localhost:8080/frontend/index.html")
    print("2. Make sure backend services are running on ports 8083, 8084")
    print("3. Test the auction platform!")
    print("")
    print("Press Ctrl+C to stop the server")
    print("==========================================")

    # Check if frontend exists
    frontend_path = os.path.join(project_dir, 'frontend', 'index.html')
    if not os.path.exists(frontend_path):
        print(f"❌ Frontend not found at: {frontend_path}")
        print("Please ensure the frontend/index.html file exists")
        sys.exit(1)

    try:
        # Create and start server
        with socketserver.TCPServer((HOST, PORT), CustomHandler) as httpd:
            print(f"✅ Server started successfully!")
            print(f"🌐 Access the frontend at: http://localhost:{PORT}/frontend/index.html")
            print("")
            httpd.serve_forever()
    except KeyboardInterrupt:
        print("\n🛑 Server stopped by user")
    except Exception as e:
        print(f"❌ Server error: {e}")
        sys.exit(1)

if __name__ == '__main__':
    main()