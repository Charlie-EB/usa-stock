const fs = require('fs');
const http = require('http');
const SftpClient = require('ssh2-sftp-client');
const { execSync } = require('child_process');

const sftp = new SftpClient();


const PRIVATE_KEY_PATH = './id_rsa';
const PUBLIC_KEY_PATH = './id_rsa.pub';

// Generate keypair if missing
if (!fs.existsSync(PRIVATE_KEY_PATH) || !fs.existsSync(PUBLIC_KEY_PATH)) {
  console.log('🔑 Generating new SSH keypair...');
  execSync(`ssh-keygen -t rsa -b 2048 -N "" -f ${PRIVATE_KEY_PATH}`);
  console.log('✅ Keypair created: id_rsa / id_rsa.pub');
}

const privateKey = fs.readFileSync(PRIVATE_KEY_PATH);
const publicKey = fs.readFileSync(PUBLIC_KEY_PATH, 'utf8');


// --- CONFIGURATION ---
const SFTP_CONFIG = {
  host: '127.0.0.1',     // your existing SFTP server
  port: 2022,            // custom SFTP port
  username: 'testuser',  // your SFTP username
  privateKey: fs.readFileSync('./id_rsa'), // path to your private key
  // passphrase: 'optional-if-your-key-is-encrypted'
};

// file to download
// const REMOTE_FILE = 'test.txt';
// const LOCAL_FILE = './downloaded_test.txt';

const REMOTE_FILE = 'sanmar_shopify.csv';
const LOCAL_FILE = './downloaded_san.csv';

// test the connection and download a file
async function testSftpDownload() {
  try {
    console.log('🔗 Connecting to SFTP server...');
    await sftp.connect(SFTP_CONFIG);

     // Debug: list files
    const list = await sftp.list('/');
    console.log('📁 Files available:', list);

    console.log(`⬇️  Downloading ${REMOTE_FILE} ...`);
    await sftp.fastGet(REMOTE_FILE, LOCAL_FILE);

    console.log(`✅ Download complete: ${LOCAL_FILE}`);
    await sftp.end();

    return { success: true, file: LOCAL_FILE };
  } catch (err) {
    console.error('❌ Error:', err.message);
    throw err;
  }
}

// optional: expose an HTTP endpoint for triggering the test
const server = http.createServer(async (req, res) => {
  if (req.url === '/test') {
    try {
      const result = await testSftpDownload();
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify(result, null, 2));
    } catch (err) {
      res.writeHead(500, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ success: false, error: err.message }));
    }
  } else {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Visit /test to connect to SFTP and download a file.');
  }
});

const PORT = 3000;
server.listen(PORT, () => {
  console.log(`🌐 HTTP test server at http://localhost:${PORT}`);
});
