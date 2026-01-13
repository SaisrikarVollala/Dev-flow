/cmd            → application entry points (main packages)
/internal       → private application code
  /api          → API logic
  /handlers     → HTTP handlers
  /services     → business logic
  /repositories → DB access
  /workers      → background jobs
  /models       → domain models
/pkg            → optional reusable libraries

