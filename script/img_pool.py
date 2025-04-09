from http.server import HTTPServer, BaseHTTPRequestHandler
import cgi
import os
import uuid
from datetime import datetime
import mimetypes
import json
import urllib.parse

# 配置
PORT = 8000
UPLOAD_DIR = "uploads"
SERVER_URL = f"http://localhost:{PORT}"  # 替换为你的服务器地址

# 确保上传目录存在
os.makedirs(UPLOAD_DIR, exist_ok=True)

# 生成唯一文件名
def generate_unique_filename(original_filename):
    ext = os.path.splitext(original_filename)[1].lower()
    date_prefix = datetime.now().strftime("%Y%m%d")
    unique_id = uuid.uuid4().hex[:8]
    return f"{date_prefix}_{unique_id}{ext}"

# 验证文件类型
def is_valid_image(filename):
    valid_extensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp']
    ext = os.path.splitext(filename)[1].lower()
    return ext in valid_extensions

class ImageServerHandler(BaseHTTPRequestHandler):
    def _set_headers(self, content_type="text/html", status_code=200):
        self.send_response(status_code)
        self.send_header('Content-type', content_type)
        # 允许跨域请求
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        self.end_headers()

    def do_OPTIONS(self):
        self._set_headers()

    def do_GET(self):
        # 主页
        if self.path == '/':
            self._set_headers()
            self.wfile.write(b"Image Pool Server is running")
            return

        # 提取文件名
        path = urllib.parse.unquote(self.path)
        if path.startswith('/images/'):
            filename = path[8:]  # 去除 '/images/' 前缀
            file_path = os.path.join(UPLOAD_DIR, filename)

            if not os.path.exists(file_path):
                self._set_headers(status_code=404)
                self.wfile.write(b"File not found")
                return

            # 获取文件类型
            content_type, _ = mimetypes.guess_type(file_path)
            if not content_type:
                content_type = 'application/octet-stream'

            # 发送文件
            with open(file_path, 'rb') as file:
                self._set_headers(content_type=content_type)
                self.wfile.write(file.read())
        else:
            self._set_headers(status_code=404)
            self.wfile.write(b"Not found")

    def do_POST(self):
        if self.path == '/upload':
            content_type = self.headers['Content-Type']
            if not content_type or not content_type.startswith('multipart/form-data'):
                self._set_headers(status_code=400)
                self.wfile.write(b"Bad request")
                return

            # 解析表单数据
            form_data = cgi.FieldStorage(
                fp=self.rfile,
                headers=self.headers,
                environ={'REQUEST_METHOD': 'POST'}
            )

            # 获取上传的文件
            if 'file' not in form_data:
                self._set_headers(status_code=400)
                self.wfile.write(b"No file part")
                return

            file_item = form_data['file']
            if not file_item.filename:
                self._set_headers(status_code=400)
                self.wfile.write(b"No selected file")
                return

            if not is_valid_image(file_item.filename):
                self._set_headers(status_code=400)
                self.wfile.write(b"Invalid file type")
                return

            # 生成唯一文件名
            unique_filename = generate_unique_filename(file_item.filename)
            file_path = os.path.join(UPLOAD_DIR, unique_filename)

            # 保存文件
            with open(file_path, 'wb') as f:
                f.write(file_item.file.read())

            # 生成URL
            file_url = f"{SERVER_URL}/images/{unique_filename}"

            # 返回成功响应
            response = {
                "filename": unique_filename,
                "url": file_url
            }
            self._set_headers(content_type="application/json")
            self.wfile.write(json.dumps(response).encode())
        else:
            self._set_headers(status_code=404)
            self.wfile.write(b"Not found")

def run(server_class=HTTPServer, handler_class=ImageServerHandler, port=PORT):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print(f"Starting image server on port {port}...")
    print(f"Upload URL: http://localhost:{port}/upload")
    print(f"Images URL: http://localhost:{port}/images/[filename]")
    httpd.serve_forever()

if __name__ == "__main__":
    run()