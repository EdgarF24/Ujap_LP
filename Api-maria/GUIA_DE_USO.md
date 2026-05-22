# Guía Rápida de Uso (Interfaz Visual - Swagger UI)

Esta guía te enseñará cómo usar tu arquitectura de microservicios utilizando tu **navegador web**, sin necesidad de abrir la terminal de comandos ni instalar herramientas externas.

Tus tres centros de control son:
1. 🔐 **Autenticación:** `http://localhost:8001/docs`
2. 🏢 **Espacios:** `http://localhost:8002/swagger/index.html`
3. 📅 **Reservas:** `http://localhost:8003/swagger-ui/index.html`

---

## PASO 1: Registrar un Administrador
*(Dónde: 🔐 Panel de Autenticación)*

1. Abre `http://localhost:8001/docs`.
2. Haz clic en la caja verde **`POST /register`**.
3. Haz clic en el botón **"Try it out"**.
4. En el campo *Request body*, pega lo siguiente:
   ```json
   {
     "email": "admin@coworking.com",
     "password": "admin123",
     "full_name": "Super Administrador",
     "role": "ADMIN"
   }
   ```
5. Haz clic en **"Execute"**. Deberías ver un código `201` en la respuesta abajo.

---

## PASO 2: Iniciar Sesión (Obtener Llave Maestra)
*(Dónde: 🔐 Panel de Autenticación)*

1. En el mismo panel, abre la caja **`POST /login`**.
2. Dale a **"Try it out"**.
3. En el formulario, escribe:
   - **username:** `admin@coworking.com`
   - **password:** `admin123`
4. Dale a **"Execute"**. 
5. En la respuesta (Código `200`), verás un texto larguísimo en el campo `"access_token"`. **Copia ese texto, es tu llave (Token).**

---

## PASO 3: Autorizar los otros Paneles
*(Dónde: 🏢 Espacios y 📅 Reservas)*

Como el sistema es seguro, para crear cosas debes decirle a Swagger que tienes permisos.
1. Abre el panel de **Espacios** (`http://localhost:8002/swagger/index.html`).
2. Haz clic en el botón verde **"Authorize"** (arriba a la derecha).
3. Pega el Token que copiaste en el Paso 2 en el recuadro que dice *Value*.
4. Dale clic a **"Authorize"** y luego a **"Close"**.
*(Repite este paso 3 en el panel de Reservas cuando vayas a reservar).*

---

## PASO 4: Crear un Espacio de Coworking
*(Dónde: 🏢 Panel de Espacios)*

1. Abre la caja verde **`POST /spaces`**.
2. Dale a **"Try it out"**.
3. En el campo *Request body*, pega lo siguiente:
   ```json
   {
     "name": "Oficina Privada 101",
     "type": "OFFICE",
     "capacity": 4,
     "price_per_hour": 25.00,
     "amenities": "WiFi, Proyector, Pizarra",
     "is_active": true
   }
   ```
4. Dale a **"Execute"**. Recibirás un `201` y el sistema te devolverá un **ID del espacio** (ej. `a3f8b...`). Copia ese ID para el próximo paso.

---

## PASO 5: Crear una Reserva
*(Dónde: 📅 Panel de Reservas)*

*Asegúrate de haber hecho el Paso 3 en este panel también.*
1. Abre `http://localhost:8003/swagger-ui/index.html`.
2. Abre la caja verde **`POST /reservations`**.
3. Dale a **"Try it out"**.
4. En el campo *Request body*, pega lo siguiente (reemplaza el `spaceId` con el que copiaste en el paso anterior):
   ```json
   {
     "spaceId": "PEGA_AQUI_EL_ID_DEL_ESPACIO",
     "startTime": "2024-12-01T10:00:00Z",
     "endTime": "2024-12-01T14:00:00Z"
   }
   ```
5. Dale a **"Execute"**. Recibirás un `201` indicando que la reserva se creó exitosamente. 

---

### ¡Listo!
Siguiendo este ciclo puedes simular a un usuario normal:
- Vas al Auth Panel, registras un usuario con `"role": "USER"`.
- Haces login con ese usuario.
- Vas al panel de Espacios -> `GET /spaces/available` para ver qué está libre.
- Vas al panel de Reservas y reservas el espacio.
