/cmd            → application entry points (main packages)
/internal       → private application code
  /api          → API logic
  /handlers     → HTTP handlers
  /services     → business logic
  /repositories → DB access
  /workers      → background jobs
  /models       → domain models
/pkg            → optional reusable libraries

Think of a folder as one unit of code.

All files in the folder:

Share variables, functions, structs

Don’t need imports between themselves

Compile together