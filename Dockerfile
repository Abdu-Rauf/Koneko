FROM ubuntu:22.04

# Install Chrome dependencies
RUN apt-get update && \
    apt-get install -y wget xvfb gnupg dbus && \
    wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | gpg --dearmor -o /usr/share/keyrings/googlechrome-linux-keyring.gpg && \
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/googlechrome-linux-keyring.gpg] http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list && \
    apt-get update && \
    apt-get install -y google-chrome-stable && \
    rm -rf /var/lib/apt/lists/*

# Create X11 socket dir
RUN mkdir -p /tmp/.X11-unix && chmod 1777 /tmp/.X11-unix

# Create non-root user
RUN useradd -m chromeuser
USER chromeuser
WORKDIR /home/chromeuser

# Run Chrome with Xvfb
CMD ["sh", "-c", "Xvfb :99 -screen 0 1280x720x24 & sleep 2 && DISPLAY=:99 google-chrome --no-sandbox --disable-gpu"]
