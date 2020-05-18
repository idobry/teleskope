<template>
    <div>
     <router-link :to="`/namespaces/${$route.params.namespace}`" style="text-decoration:none; color:inherit;" >
    <page-title :heading=$route.params.namespace :subheading=subheading :icon=icon></page-title>
     </router-link>
             <div class="row mt-4">
            <div class="col-lg-6 col-xl-6"  v-for="(chartValues,deployment,index) in deployments" :key="index" >    
              <router-link :to="`/namespaces/${$route.params.namespace}/deployments/${deployment}`" style="text-decoration:none; color:inherit;" >
                <div class="card mb-3 widget-content">
                    <div class="widget-content-wrapper">
                        <div class="widget-content-left">
                            <div class="widget-superheading"><span>{{ deployment }}</span></div>
                        </div>
                        <div class="widget-content-right">
                            <doughnut :chartValues="chartValues" style="width:50px; height:50px"></doughnut>
                        </div>
                    </div>
                </div>
              </router-link>
            </div>
        </div>
    </div>
</template>

<script>

    import PageTitle from "../../Layout/Components/PageTitle.vue";
    import ApiService from '../../Services/apiService'
    
import doughnut from '../Charts/Chartjs/Doughnut'

    export default {
        components: {
            PageTitle,
            doughnut,
        },
        data: () => ({
            apiService: new ApiService(),
            heading: 'Chart Boxes III',
            subheading: null,
            icon: 'pe-7s-helm icon-gradient bg-plum-plate',
            deployments: {},
        }),
        async beforeMount() {
            this.apiService.connectToWebSocket(this.onmessage);
            const deploymentsList = await this.apiService.listDeploymentsAt(this.$route.params.namespace)
            deploymentsList.forEach(dep => {
                this.deployments[dep] = {
                      exist: 1,
                      toBeCreated: 0,
                }
            })
            this.subheading = `${Object.keys(this.deployments).length || 0 } Deployments:`
        },
        beforeDestroyed() {
            this.apiService.disconnectFromWebSocket();
      },
        methods: {
            onmessage(msg) {
                if(msg.Namespace === this.$route.params.namespace && this.deployments[msg.Name]) {
                    const newChartvalues = {
                        exist: msg.ReplicaCurrent,
                        toBeCreated: Math.max(msg.ReplicaDesired - msg.ReplicaCurrent,0),
                    }
                    this.deployments[msg.Name] = newChartvalues;
                    this.$forceUpdate()
                }
            }
        }
    }

</script>
