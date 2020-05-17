import axios from 'axios'

export default class ApiService {
    
    constructor(){
        // this.apiBaseURL = process.env['VUE_APP_API_BASE_URL'] || 'http://localhost:3000' 
        this.axios = axios.create({
            // baseURL: this.apiBaseURL,
            /* other custom settings */
        });
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

}
