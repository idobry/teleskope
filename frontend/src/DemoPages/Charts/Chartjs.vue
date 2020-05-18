<template>
  <div>
     <router-link :to="`/namespaces/${$route.params.namespace}`" style="text-decoration:none; color:inherit;" >
    <page-title :heading=$route.params.deployment :subheading=subheading :icon=icon></page-title>
     </router-link>
    <div class="content">
        <b-row>
          <b-col  md="6">
            <b-card title="Running Pods:" v-if="chartValues" class="main-card mb-3">
              <doughnut :chartValues="chartValues" ></doughnut>
            </b-card>
          </b-col>
            <div class="col-lg-6 col-xl-6">
                <div class="card mb-3 widget-content">
                    <div class="widget-content-wrapper">
                        <div class="widget-content-left">
                            <div class="widget-heading text-success">Image</div>
                            <div class="widget-heading">{{ image }}</div>
                        </div>
                        <div class="widget-content-right">
                        </div>
                    </div>
                </div>

        <div class="divider mt-0" style="margin-bottom: 17px;"></div>
                <div class="card mb-3 widget-content">
                    <div class="widget-content-wrapper">
                        <div class="widget-content-left">
                            <div class="widget-heading text-success">Tag</div>
                            <div class="widget-heading">{{ tag }}</div>
                        </div>
                        <div class="widget-content-right">
                        </div>
                    </div>
                </div>

        <div class="divider mt-0" style="margin-bottom: 17px;"></div>
                <div class="card mb-3 widget-content">
                    <div class="widget-content-wrapper">
                        <div class="widget-content-left">
                            <div class="widget-heading">Age</div>
                        </div>
                        <div class="widget-content-right">
                            <div class="widget-numbers text-success"><span>2m31s</span></div>
                        </div>
                    </div>
                </div>
                <div class="divider mt-0" style="margin-bottom: 17px;"></div>
                  <div class="card mb-3 widget-content">
                    <div class="widget-content-wrapper">
                        <div class="widget-content-left">
                            <div class="widget-heading">Pods Count</div>
                        </div>
                        <div class="widget-content-right" v-if="podCount">
                            <div class="widget-numbers text-success"><span> {{ podCount.current }} / {{ podCount.desired }} </span></div>
                        </div>
                    </div>
                </div>
            </div>
        </b-row>
    <b-card v-if="envVars" title="Environment Variables" class="main-card mb-4">
        <b-table :striped="true"
                 :bordered="false"
                 :outlined="false"
                 :small="false"
                 :hover="true"
                 :dark="false"
                 :fixed="false"
                 :foot-clone="false"
                 :items="envVars"
                 :fields="fields">
        </b-table>
    </b-card>
    </div>
  </div>
</template>

<script>

  import PageTitle from "../../Layout/Components/PageTitle.vue";

  import doughnut from './Chartjs/Doughnut'
  import radar from './Chartjs/Radar'
  import polar from './Chartjs/Polar'
  import pie from './Chartjs/Pie'
  import lineeg from './Chartjs/Line'
  import areaeg from './Chartjs/Area'
  import bar from './Chartjs/Bar'
  import barhoriz from './Chartjs/BarHoriz'
  import ApiService from '../../Services/apiService'

  export default {
    components: {
      PageTitle,

      doughnut,
      radar,
      polar,
      pie,
      lineeg,
      areaeg,
      bar,
      barhoriz,

    },
    data: () => ({
      apiService: new ApiService(),
      ws: null,
      subheading: null,
      image: null,
      tag: null,
      age: null,
      chartValues: null,
      fields: [ 'name', 'value' ],
      envVars: null,
      podCount: null,
      icon: 'pe-7s-helm icon-gradient bg-amy-crisp',
    }),
    async mounted() {
        this.apiService.connectToWebSocket(this.onmessage);
        await this.setInitialFeed();
      },
    beforeDestroyed() {
        this.apiService.disconnectFromWebSocket();
      },
    computed: {
    },
    methods: {
      getEnvVars(container) {
        return container && container.Envs.map(pair => {
          const [name, ...valueArray] = pair.split(':')
          return {
            name,
            value: valueArray.join(':')
          }
        })
      },
      getImage(container) {
        return container && container.Image.split(':')[0]
      },
      getTag(container) {
        return container && container.Image.split(':')[1]
      },
      getAge(container) {
        return '2m31s'
      },
      getSubheading(container) {
        return 'Current State: success'
      },
      getChartValues(msg) {
        return msg && {
          exist: msg.ReplicaCurrent,
          toBeCreated: Math.max(msg.ReplicaDesired - msg.ReplicaCurrent,0),
          toBeDeleted: Math.max(msg.ReplicaCurrent - msg.ReplicaDesired,0)
        }
      },
      getPodCount(msg) {
        return msg && {
          current: msg.ReplicaCurrent,
          desired: msg.ReplicaDesired,
        }
      },
      getMasterContainer(msg) {
        return msg && msg.Containers && msg.Containers.find(container => container.Name == 'master')
      },
      async setInitialFeed(){
        const message = await this.apiService.describeDeployment(this.$route.params.namespace,this.$route.params.deployment)
        this.updateFeed(message)
      },
      updateFeed(message) {
        console.log(message);
          this.chartValues = this.getChartValues(message);
          this.podCount = this.getPodCount(message);
          const container = this.getMasterContainer(message);
          this.image = this.getImage(container)
          this.tag = this.getTag(container)
          this.age = this.getAge(container)
          this.subheading = this.getSubheading(container)
          this.envVars = this.getEnvVars(container)
      },
      onmessage(message) {
          message.Name === this.$route.params.deployment && this.updateFeed(message)
      },
    }
}
</script>
