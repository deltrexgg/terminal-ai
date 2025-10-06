

  <h1>Terminal AI — Local Code Debugging Assistant</h1>

  <p class="muted">A small CLI tool (<code>tai</code>) that sends code or questions to a locally hosted Qwen GGUF model running in server mode, and prints a clean assistant reply in the terminal.</p>

  <h2>Overview</h2>
  <p>This project includes two parts:</p>
  <ol>
    <li>Local AI model server based on a Qwen <code>.gguf</code> model and a runtime <code>.llamafile</code>, exposing an OpenAI-compatible HTTP API.</li>
    <li>A Go-based CLI tool (<code>tai</code>) that posts requests to that API and prints only the assistant's message.</li>
  </ol>

  <h2>Prerequisites</h2>
  <ul>
    <li>Linux or macOS machine for the CLI. For the model server: Linux recommended (or an appropriate binary for macOS).</li>
    <li>Go (for building the CLI) if you want to build from source.</li>
    <li>Qwen model files:</li>
  </ul>
  <pre><code>Qwen2.5-0.5B-Instruct-Q6_K.gguf
Qwen2.5-0.5B-Instruct-Q6_K.llamafile</code></pre>

  <h2>Run the model server (direct)</h2>
  <p>From the directory containing the two files:</p>
  <pre><code>nohup ./Qwen2.5-0.5B-Instruct-Q6_K.llamafile \
  --model Qwen2.5-0.5B-Instruct-Q6_K.gguf \
  --server \
  --port 8081 \
  --n-gpu-layers 20 \
  --threads 8 \
  --host 0.0.0.0 \
  > qwen.log 2>&1 & echo $! > qwen.pid</code></pre>
  <p>Default API endpoint: <code>http://&lt;server-ip&gt;:8081</code></p>
  <p>Test with curl:</p>
  <pre><code>curl http://127.0.0.1:8081/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"Qwen2.5-0.5B-Instruct-Q6_K","messages":[{"role":"user","content":"Hello"}]}'</code></pre>

  <h2>Run the model server (Docker)</h2>
  <p>Example Dockerfile (place in same folder as the model files):</p>
  <pre><code>FROM ubuntu:22.04
RUN apt-get update && apt-get install -y wget curl vim && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY Qwen2.5-0.5B-Instruct-Q6_K.llamafile .
COPY Qwen2.5-0.5B-Instruct-Q6_K.gguf .
RUN chmod +x Qwen2.5-0.5B-Instruct-Q6_K.llamafile
EXPOSE 8081
ENTRYPOINT ["./Qwen2.5-0.5B-Instruct-Q6_K.llamafile"]
CMD ["--model", "Qwen2.5-0.5B-Instruct-Q6_K.gguf", "--server", "--port", "8081", "--n-gpu-layers", "20", "--threads", "8", "--host", "0.0.0.0"]</code></pre>

  <p>Build and run:</p>
  <pre><code>docker build -t qwen-ai .
docker run -d --name qwen -p 8081:8081 qwen-ai</code></pre>

  <p>Notes:</p>
  <ul>
    <li>If your llamafile binary is for a different CPU architecture than the host (for example ARM vs amd64), either use a matching base image or download the binary for the correct architecture.</li>
    <li>To run with host networking (useful for quick LAN access): <code>docker run -d --network host qwen-ai</code></li>
  </ul>

  <h2>CLI: Build and install</h2>
  <p>Build the Go binary for your platform:</p>
  <pre><code>GOOS=linux GOARCH=amd64 go build -o tai main.go
# or for macOS Apple Silicon:
GOOS=darwin GOARCH=arm64 go build -o tai main.go</code></pre>
  <p>Install system-wide:</p>
  <pre><code>sudo mv tai /usr/local/bin/
sudo chmod +x /usr/local/bin/tai</code></pre>
  <p>On macOS, if you copied the binary from another device and get "operation not permitted", remove quarantine:</p>
  <pre><code>sudo xattr -d com.apple.quarantine /usr/local/bin/tai</code></pre>

  <h2>Usage</h2>
  <p>Set API URL (optional). If not set, CLI falls back to <code>http://192.168.29.200:8081</code> in example code.</p>
  <pre><code>export TAIURL="http://192.168.29.200:8081"</code></pre>

  <p>Run CLI:</p>
  <pre><code>tai "Here is my Go function:\n\nfunc add(a, b int) int { return a + b }" 300</code></pre>

  <p>Command format:</p>
  <table>
    <tr><th>Argument</th><th>Description</th><th>Default</th></tr>
    <tr><td>question</td><td>Code or text to send to the assistant</td><td>required</td></tr>
    <tr><td>max_tokens</td><td>Maximum tokens to request from the model</td><td>200</td></tr>
  </table>

  <h2>Example request body (sent by the CLI)</h2>
  <pre><code>{
  "model": "Qwen2.5-0.5B-Instruct-Q6_K",
  "messages": [
    {"role": "system", "content": "You are an expert code debugging assistant. Always explain errors, suggest fixes, and provide corrected code."},
    {"role": "user", "content": "/* code or question */"}
  ],
  "max_tokens": 200,
  "temperature": 0
}</code></pre>

  <h2>Parsing the response</h2>
  <p>The server returns an OpenAI-compatible JSON object. The CLI unmarshals it and prints only the assistant text located at <code>choices[0].message.content</code>.</p>

  <h2>Packaging as a .deb (optional)</h2>
  <ol>
    <li>Create layout: <code>terminal-ai_1.0.0/usr/local/bin/terminal-ai</code></li>
    <li>Create control file at <code>terminal-ai_1.0.0/DEBIAN/control</code> with package metadata.</li>
    <li>Build: <code>dpkg-deb --build terminal-ai_1.0.0</code></li>
    <li>Install: <code>sudo apt install ./terminal-ai_1.0.0.deb</code></li>
  </ol>

  <h2>Troubleshooting</h2>
  <h3>Connection refused</h3>
  <p>Confirm the server is running and listening on the expected port:</p>
  <pre><code>curl http://127.0.0.1:8081/v1/models</code></pre>

  <h3>Container exits immediately with "exec format error"</h3>
  <p>Binary architecture mismatch. On the host run <code>uname -m</code>. Ensure the llamafile is built for that architecture or build an image with the matching platform (for example <code>FROM --platform=linux/arm64 ubuntu:22.04</code>).</p>

  <h3>macOS "operation not permitted"</h3>
  <pre><code>sudo xattr -d com.apple.quarantine /usr/local/bin/tai
sudo chmod +x /usr/local/bin/tai</code></pre>

  <h2>Security</h2>
  <ul>
    <li>If exposing the API over a network, restrict access via firewalls or run behind an authenticated proxy.</li>
    <li>Do not expose the model server publicly without proper access control.</li>
  </ul>

  <h2>Future improvements</h2>
  <ul>
    <li>Interactive REPL mode for CLI</li>
    <li>Streaming responses from the server to show incremental output</li>
    <li>Config file for defaults (e.g., <code>~/.tairc</code>)</li>
    <li>Systemd unit for automatic restart on server hosts</li>
  </ul>

  <h2>License</h2>
  <p>MIT</p>

  <hr />

