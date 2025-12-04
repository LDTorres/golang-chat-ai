# Chat RAG App — Go + Fiber + Pulumi + RAG

Cgat AIT con Golang y LLM + RAG

## Especificaciones:
	•	Go + Fiber → Backend web simple con endpoints de chat
	•	Frontend minimalista → Página web con interfaz de chat
	•	RAG (Retrieval-Augmented Generation) → Conexión con un Vector DB y un LLM
	•	Infraestructura en la nube con Pulumi → AWS (ECS Fargate, RDS, S3, etc.)

El objetivo es construir un sistema moderno que incluya backend, frontend, IA, vector search y IaC.

⸻

🚀 Objetivos del proyecto
	•	Crear un chat básico que envía un mensaje al backend y recibe una respuesta generada por IA.
	•	Integrar un pipeline de RAG para que la IA use información externa almacenada en una base vectorial.
	•	Desplegar toda la infraestructura con Pulumi en AWS (infra reproducible).
	•	Arquitectura cloud, contenedores, redes y buenas prácticas.

⸻

🧩 Arquitectura General

Componentes principales:
	•	Frontend: página HTML/JS simple servida por Fiber.
	•	Backend: API Go + Fiber con:
	•	Manejo de chat
	•	Orquestación de RAG
	•	Conexión al LLM
	•	Conexión al Vector DB
	•	RAG:
	•	Vector DB (Qdrant Cloud o Pinecone)
	•	Embeddings + Retrieval
	•	Prompting con contexto
	•	Infra:
	•	AWS ECS Fargate (app)
	•	ALB (load balancer)
	•	S3 (documentos a indexar)
	•	RDS Postgres (historial de conversaciones)
	•	VPC, subnets, SG, IAM roles
	•	SSM/Secrets Manager

⸻

🏗️ Infraestructura (Pulumi)

Pulumi despliega:
	•	VPC + subnets públicas/privadas
	•	Security groups
	•	S3 bucket para documentos RAG
	•	RDS Postgres (opcional)
	•	ECR repository
	•	ECS Fargate Cluster
	•	ECS Service para la app de Go
	•	Application Load Balancer
	•	SSM/Secrets Manager para:
	•	API keys de LLM
	•	API keys del vector DB
	•	Credenciales de DB

Puedes tener 1 stack por entorno:

/infra
  Pulumi.dev.yaml
  Pulumi.prod.yaml
  main.go (o index.ts)


⸻

🔧 Requisitos

Local
	•	Go 1.22+
	•	Docker / Docker Compose
	•	Pulumi CLI
	•	AWS CLI configurado
	•	Access tokens del LLM provider
	•	Access tokens del vector DB

Servicios externos
	•	Qdrant Cloud / Pinecone
	•	OpenAI / Groq / Anthropic / (el que quieras)

⸻

🧪 Ejecución local

1. Clonar repo

2. Variables de entorno

Crear .env:

OPENAI_API_KEY=...
VECTOR_DB_URL=...
VECTOR_DB_API_KEY=...
DATABASE_URL=postgres://...
S3_BUCKET_NAME=...

3. Correr local con Docker Compose

docker compose up --build

App disponible en:

👉 http://localhost:8080

⸻

🧠 Flujo RAG (simplificado)
	1.	Usuario envía mensaje → /api/chat
	2.	Backend genera embedding del mensaje
	3.	Busca contexto en el vector DB
	4.	Construye prompt → contexto + pregunta
	5.	Envía a LLM
	6.	Devuelve respuesta + contexto usado

⸻

📥 Ingesta de documentos

Para indexar documentos en el vector DB:

go run cmd/ingest/main.go --path=./docs

El pipeline:
	•	Carga archivos desde docs/
	•	Chunking (división en fragmentos)
	•	Generación de embeddings
	•	Inserción en el vector DB

⸻

🧱 Estructura del repositorio

Ejemplo sugerido:

/cmd
  /app
    main.go
  /ingest
    main.go

/internal
  /http
    handler.go
    router.go
  /service
    chat.go
    rag.go
  /llm
    client.go
  /vector
    qdrant.go
  /storage
    s3.go
  /db
    postgres.go

/frontend
  index.html
  styles.css
  app.js

/infra
  main.go (o index.ts si usas TS para Pulumi)

/docs
  (documentos para RAG)

Dockerfile
docker-compose.yaml
README.md


⸻

☁️ Despliegue en AWS con Pulumi

1. Configurar stack

pulumi stack select dev
pulumi config set aws:region us-east-1

2. Deploy completo

pulumi up

Pulumi creará:
	•	VPC + subnets
	•	ALB
	•	ECS + Fargate
	•	Postgres
	•	S3
	•	Secrets

Endpoints mostrados al final del pulumi up.

⸻

🎯 Roadmap del Proyecto

Fase 1 — Backend + Frontend local

Fase 2 — Integración con LLM

Fase 3 — RAG local (Qdrant + Docker)

Fase 4 — Infra Pulumi (ECS sin RAG)

Fase 5 — RAG en la nube

Fase 6 — WebSocket + Streaming

Fase 7 — Login de usuarios (opcional)

Fase 8 — UI más completa (opcional)

⸻

🤝 Contribuciones

Pull Requests bienvenidos.
Issues también.

⸻

📜 Licencia

MIT (o la que prefieras).

⸻

Si querés, puedo generar también:

✅ Un diagrama PNG para agregar al README
✅ Un CONTRIBUTING.md
✅ Un Makefile (build, run, deploy)
✅ Plantilla inicial del repositorio (carpetas + archivos vacíos)

¿Querés que lo prepare?