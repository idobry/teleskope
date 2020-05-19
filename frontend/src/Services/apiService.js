import axios from 'axios'

export default class ApiService {
    
    constructor(){
        this.apiBaseURL = process.env['VUE_APP_API_BASE_URL'] || 'http://localhost:3000' 
        this.axios = axios.create({
            baseURL: this.apiBaseURL,
            /* other custom settings */
        });
        this.webSocket = null;
        this.wsEndpoint = process.env.VUE_APP_WS_ENDPOINT || 'ws://localhost:3000/ws';
    }

    async getAPI(route) {
        console.log(`calling ${route}:`)
        try {
            const response = await this.axios.get(route)
            return response.data
        } catch (error) {
            console.log(`error getting ${route}:`)
            console.log(error); 
      }
    }

    async listNamespace() {
        const list = await this.getAPI('/list/ns')
        return list && list.ID
    }

    async listDeploymentsAt(namespace) {
        const list = await this.getAPI(`/list/dep/${namespace}`)
        return list && list.ID
    }

    async describeDeployment(namespace,deployment) {
        return await this.getAPI(`/dep/${namespace}/${deployment}`)
    }

    connectToWebSocket(onmessage = () => {}) {
        if (this.webSocket !== null) {
            this.webSocket.close()
        }

        this.webSocket = new WebSocket(this.wsEndpoint);

        this.webSocket.onmessage = evt =>  {
            console.log("ommessage.");
            const message = JSON.parse(evt.data);
            console.log(message);
            onmessage(message)
        };

        this.webSocket.onclose = evt =>  {
          console.log("onclose.");
            this.webSocket = new WebSocket(this.wsEndpoint);
        };

        this.webSocket.onopen = evt =>  {
          console.log("onopen.");
        };

        this.webSocket.onerror = evt =>  {
            console.log("Error!");
        };
  }
  disconnectFromWebSocket() {
      this.webSocket.close();
  }

}
