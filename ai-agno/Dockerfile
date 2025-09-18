FROM python:3.11-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

# Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /usr/local/bin/uv

# Set working directory
WORKDIR /app/ai-agno

# Copy requirements first for better caching
COPY requirements.txt* ./

# Create virtual environment and install dependencies with uv
RUN uv venv .venv && \
    uv pip install --python .venv -r requirements.txt

# Copy application code
COPY . .

# Activate virtual environment for runtime
ENV PATH="/app/ai-agno/.venv/bin:${PATH}"
ENV VIRTUAL_ENV="/app/ai-agno/.venv"

# Run the main application
CMD ["python", "main.py"]