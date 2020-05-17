<template>
    <div class="app-page-title">
        <div class="page-title-wrapper">
            <div class="page-title-heading">
                <div class="page-title-icon">
                    <i :class="icon"/>
                </div>
                <div>
                    {{heading}}
                    <div
                        class="page-title-subheading">
                        {{subheading}}
                    </div>
                </div>
            </div>
            <div class="page-title-actions">
                <button type="button" class="btn-shadow mr-3 btn btn-dark" @click="reloadPage()">
                    <font-awesome-icon icon="redo"/>
                </button>
                 <router-link :to="`/namespaces/${$route.params.namespace}`" style="text-decoration:none; color:inherit;" v-if="$route.name == 'deployment'">
                <button type="button" class="btn-shadow d-inline-flex align-items-center btn btn-success">
                    <font-awesome-icon class="mr-2" icon="angle-left"/>
                    Back to {{ $route.params.namespace }} namespace
                </button>
                 </router-link>
            </div>
        </div>
    </div>
</template>

<script>
    import {library} from '@fortawesome/fontawesome-svg-core'
    import {
        faStar,
        faPlus,
        faRedo,
        faAngleLeft
    } from '@fortawesome/free-solid-svg-icons'
    import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome'

    library.add(
        faStar,
        faPlus,
        faRedo,
        faAngleLeft
    );

    export default {
        components: {
            'font-awesome-icon': FontAwesomeIcon,
        },
        props: {
            icon: String,
            heading: String,
            subheading: String,

        },
        methods: {
            reloadPage() {
                window.location.reload()
            }
        }
    }
</script>
