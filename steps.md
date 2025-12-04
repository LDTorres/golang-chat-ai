1. Alcance del sistema

Funcionalidad inicial:
	•	Página web sencilla con UI de chat.
	•	Backend en Go + Fiber (HTTP + WebSocket/SSE para streaming). Fiber está pensado justo para este tipo de APIs rápidas y soporta bien websockets y middleware.
	•	El chat:
	•	Envía mensaje del usuario.
	•	Llama a la AI usando RAG.
	•	Devuelve la respuesta al frontend.
	•	Infra en la nube levantada y destruida con Pulumi (idealmente AWS + ECS Fargate).

⸻

2. Arquitectura lógica (vista alto nivel)

Imagina algo así:

Cliente (browser)
⬇️⬆️ HTTP/WebSocket
Fiber App (Chat + RAG API)
⬇️
	•	LLM API (por ej. OpenAI / otro provider)
	•	Vector DB (Qdrant / Pinecone / etc.)
	•	Storage de documentos (S3 u otro)
	•	DB relacional opcional (Postgres para historial de chat / usuarios)

2.1. Frontend web

Súper simple:
	•	Una página servida por Fiber:
	•	Campo de texto.
	•	Lista de mensajes.
	•	Botón "Enviar".
	•	Comunicación:
	•	Primera versión: HTTP POST /api/chat.
	•	Versión posterior: WebSocket o SSE para ver el texto de la IA "streaming".

2.2. Servicio backend en Go + Fiber

Responsabilidades:
	1.	Endpoints HTTP/WS
	•	GET / → sirve el HTML/JS/CSS.
	•	POST /api/chat → recibe {message, conversationId} y devuelve {response, sources}.
	•	Opcional: /api/docs/upload para subir documentos al RAG.
	2.	Capa de aplicación (ChatService)
	•	Valida input.
	•	Llama a RagOrchestrator.
	•	Persiste historial de mensajes (si usas DB).
	3.	RagOrchestrator
	•	Recibe userMessage.
	•	Llama a Retriever (vector DB).
	•	Construye el prompt con contexto.
	•	Llama al LLMClient.
	•	Devuelve la respuesta + chunks de contexto usados.
	4.	Integraciones
	•	EmbeddingClient → provider de embeddings.
	•	LLMClient → provider de texto.

⸻

3. Diseño del RAG dentro del backend

3.1. Componentes de RAG
	•	Document Store:
	•	Raw docs (PDFs/markdown) en S3 o similar.
	•	Vector Store:
	•	Para práctica técnica: Qdrant (open source, muy usado para RAG).
	•	Puedes correrlo en Docker (local) y luego como contenedor en cloud.
	•	Ingestion Pipeline (puede ser un comando CLI o endpoint admin):
	1.	Leer documentos (desde un folder o S3).
	2.	Partirlos en chunks.
	3.	Pedir embeddings al provider (OpenAI u otro).
	4.	Guardar (id, texto, metadata, vector) en Qdrant.
	•	Query Pipeline (en cada mensaje):
	1.	Generar embedding del mensaje.
	2.	Buscar vecinos en Qdrant (top-k, filtros, etc.).
	3.	Armar prompt: "User pregunta X; contexto: [docs]".
	4.	Llamar al LLM.
	5.	Devolver respuesta + contexto (para debug en UI si quieres).

⸻

4. Persistencia

Depende qué quieres practicar, pero algo razonable:
	•	Vector DB: Qdrant (container) o servicio gestionado.
	•	Relacional (opcional pero educativo):
	•	Postgres (RDS o Aurora Serverless v2):
	•	users
	•	conversations
	•	messages (quién envió, timestamp, texto, si vino del modelo).
	•	Object Storage:
	•	S3 bucket para:
	•	Docs originales.
	•	Logs o dumps de embeddings si quieres.

⸻

5. Arquitectura física / despliegue con Pulumi

Asumiendo AWS (porque los ejemplos de Pulumi vienen muy armados para ECS Fargate).

5.1. Componentes de infra

Pulumi (en Go o TS, como prefieras) aprovisiona:
	1.	Red
	•	VPC
	•	Subnets públicas/privadas
	•	Security groups (HTTP/HTTPS desde internet al ALB, interno al ECS/Qdrant/DB).
	2.	ECS Fargate + Load Balancer
	•	ECR repo para tu imagen de Fiber.
	•	ECS Cluster.
	•	Task Definition + Service:
	•	Container de tu app Go + Fiber.
	•	Application Load Balancer (ALB):
	•	Listener 80/443 → target group → ECS Service.
	•	Puedes usar el template de "container service" de Pulumi para ACS Fargate como base.
	3.	Qdrant
Opciones:
	•	Sencilla para empezar: Qdrant Cloud (SaaS) y no lo manejas con Pulumi (sólo guardas la URL/API key en SSM).
	•	Más IaaC práctica: Segundo ECS service con la imagen oficial de Qdrant, accesible sólo dentro de la VPC.
	4.	Base de datos (si la usas)
	•	RDS Postgres (subnet privada).
	•	Seguridad: sólo accesible desde las tasks de ECS.
	5.	S3 Bucket
	•	Para almacenar documentos.
	•	Política para que solo tu app (rol de tarea ECS) pueda leer/escribir.
	6.	Configuración / secretos
	•	SSM Parameter Store o Secrets Manager:
	•	API key de LLM.
	•	API key de embeddings.
	•	Credenciales de Qdrant (si aplica).
	•	Variables de entorno en la Task Definition leyendo de SSM.
	7.	Pulumi stacks
	•	Pulumi.dev.yaml, Pulumi.prod.yaml.
	•	Puedes controlar:
	•	Tamaños de Fargate tasks.
	•	URLs de Qdrant (local, cloud, etc.).
	•	Clave de OpenAI de dev vs prod.

⸻

6. Flujo completo de un mensaje (end-to-end)
	1.	Usuario escribe mensaje en el browser y hace clic en "Enviar".
	2.	Frontend hace POST /api/chat con {message, conversationId}.
	3.	Fiber:
	•	Valida input.
	•	Llama a ChatService.HandleMessage.
	4.	ChatService:
	•	Si quieres, guarda el mensaje en Postgres.
	•	Llama a RagOrchestrator.Answer(message, conversationContext).
	5.	RagOrchestrator:
	•	Pide embedding al EmbeddingClient.
	•	Hace search en Qdrant por los top-k documentos relevantes.
	•	Construye el prompt del LLM con:
	•	Instrucciones del sistema.
	•	Contexto (chunks).
	•	Mensaje del usuario.
	•	Llama al LLMClient.
	6.	LLM responde.
	7.	RagOrchestrator devuelve {answer, retrievedDocs}.
	8.	ChatService:
	•	Opcional: guarda la respuesta en Postgres.
	•	Devuelve al handler.
	9.	Fiber responde al frontend.
	10.	Frontend actualiza la UI del chat.

Más adelante, puedes cambiar el paso 2-10 a WebSocket para streaming.

⸻

7. Roadmap de implementación (orden recomendado)

Fase 1 - Todo local sin RAG
	•	Fiber server + endpoint /api/chat que responde "Echo: …".
	•	UI mínima de chat.
	•	Correrlo con docker-compose.

Fase 2 - Integrar LLM "pelado"
	•	En el backend, cambiar el "Echo" por llamada al LLM (sin RAG).
	•	Dejas todo local.

Fase 3 - RAG local
	•	Correr Qdrant en Docker localmente.
	•	Crear un pequeño script/CLI de ingestion que:
	•	Lee unos cuantos .md/.txt de ejemplo,
	•	Genera embeddings,
	•	Inserta en Qdrant.
	•	Integrar el Retriever en tu endpoint /api/chat.

Fase 4 - Infra con Pulumi
	•	Tomar el template de Container Service en AWS con Pulumi (ECS Fargate + ALB).
	•	Adaptarlo para tu imagen de Fiber.
	•	Desplegar tu app en Fargate (sin Qdrant aún, sólo LLM simple).

Fase 5 - RAG en la nube
	•	Opción rápida: usar Qdrant Cloud/Pinecone y configurarlo vía env vars.
	•	Opción full-IaC: desplegar Qdrant también en ECS con Pulumi.
	•	Apuntar tu app en Fargate a ese vector DB.