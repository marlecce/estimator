# Estimator

Estimator is a collaborative task estimation application designed to streamline the process of estimating tasks within a team. It enables hosts to create estimation sessions, invite participants, collect anonymous estimates, and reveal the results in a controlled manner.

## Features

- Room Management: Hosts can create and share estimation rooms.

- Anonymous Estimation: Participants can provide estimates without revealing their input.

- Controlled Results: Only the host can reveal the submitted estimations.

- Real-Time Updates: Seamless interaction powered by WebSockets.

## Architecture

Estimator is implemented as a monorepo with the following components:

- Backend: A RESTful API built with Go.
- Frontend: A user-friendly interface developed in Vue.js.

- CI/CD: Automated pipelines for testing, building, and deploying using GitHub Actions.

## Repository Structure

```bash
.
├── backend/       # Go-based backend code
├── frontend/      # Vue.js-based frontend code
├── Makefile       # Build and run commands
├── docker-compose.yml  # Docker Compose configuration
└── README.md      # Project documentation
```

## Installation

### Prerequisites

- Go 1.23+

- Node.js 23+

- Docker and Docker Compose

## Setup

1. Clone the repository:

```bash
git clone https://github.com/yourusername/estimator.git
cd estimator
```

2. Install dependencies:

```bash
make install-frontend
```

3. Build the project:

```bash
make all
```

4. Run the development environment:

```bash
make run-backend
make dev-frontend
```

## Docker Usage

Build and run the application using Docker:

```bash
docker compose up
```
