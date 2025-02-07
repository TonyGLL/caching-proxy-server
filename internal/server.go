package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func StartServer(port int, origin string) {
	//* Iniciar el servidor proxy
	fmt.Printf("Starting caching proxy on port %d with origin %s\n", port, origin)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestedURL := r.URL.String()
		fmt.Printf("Requested URL: %s\n", requestedURL)

		//* Verificar si la respuesta está en Redis
		cachedResponse, err := rdb.Get(ctx, requestedURL).Result()
		if err == nil {
			//* Si existe en Redis, devolver la respuesta almacenada
			w.Header().Set("X-Cache", "HIT")
			var cachedData map[string]interface{}
			if err := json.Unmarshal([]byte(cachedResponse), &cachedData); err != nil {
				http.Error(w, "Error decoding cached response", http.StatusInternalServerError)
				return
			}

			//* Devolver los headers almacenados
			headers := cachedData["headers"].(map[string]interface{})
			for key, value := range headers {
				w.Header().Set(key, value.(string))
			}

			//* Devolver la respuesta almacenada
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(cachedData["response"].(string)))
			return
		}

		//* Si no está en Redis, hacer la solicitud al servidor de origen
		w.Header().Set("X-Cache", "MISS")

		targetURL := origin + requestedURL

		//* Hacer la solicitud al servidor de origen
		resp, err := http.Get(targetURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching from origin: %v", err), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		//* Verificar si la respuesta del servidor de origen es válida
		if resp.StatusCode >= 400 {
			http.Error(w, "Origin server returned an error", http.StatusBadGateway)
			return
		}

		//* Leer la respuesta del servidor de origen
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			http.Error(w, fmt.Sprintf("Error decoding origin response: %v", err), http.StatusInternalServerError)
			return
		}

		//* Convertir la respuesta a JSON
		responseJSON, err := json.Marshal(result)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response to JSON: %v", err), http.StatusInternalServerError)
			return
		}

		//* Guardar los headers de la respuesta
		headers := make(map[string]string)
		for key, values := range resp.Header {
			headers[key] = values[0] //* Tomamos el primer valor del slice
		}

		//* Crear un objeto para almacenar en Redis
		cacheData := map[string]interface{}{
			"response": string(responseJSON),
			"headers":  headers,
		}

		//* Convertir el objeto a JSON
		cacheDataJSON, err := json.Marshal(cacheData)
		if err != nil {
			fmt.Printf("Error encoding cache data to JSON: %v\n", err)
		} else {
			//* Almacenar la respuesta y los headers en Redis
			err = rdb.Set(ctx, requestedURL, cacheDataJSON, 10*time.Minute).Err() //* Expira en 10 minutos
			if err != nil {
				fmt.Printf("Error saving response to Redis: %v\n", err)
			} else {
				fmt.Printf("Response and headers saved in Redis for URL: %s\n", requestedURL)
			}
		}

		//* Devolver los headers al cliente
		for key, value := range headers {
			w.Header().Set(key, value)
		}

		//* Devolver la respuesta al cliente
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	})

	//* Iniciar el servidor
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

// * ClearCache limpia la caché
func ClearCache() error {
	//* Limpiar la caché
	fmt.Println("Clearing cache...")
	return rdb.FlushDB(ctx).Err()
}
