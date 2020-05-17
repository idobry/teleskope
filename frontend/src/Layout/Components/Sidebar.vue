<template>
    <div class="app-sidebar sidebar-shadow" @mouseover="toggleSidebarHover('add','closed-sidebar-open')" @mouseleave="toggleSidebarHover('remove','closed-sidebar-open')">
        <div class="app-header__logo">
            <div class="logo-src"/>
            <div class="header__pane ml-auto">
                <button type="button" class="hamburger close-sidebar-btn hamburger--elastic" v-bind:class="{ 'is-active' : isOpen }" @click="toggleBodyClass('closed-sidebar')">
                    <span class="hamburger-box">
                        <span class="hamburger-inner"></span>
                    </span>
                </button>
            </div>
        </div>
        <div  class="app-sidebar-content">
            <VuePerfectScrollbar v-if="menu"  class="app-sidebar-scroll" v-once>
                <sidebar-menu  showOneChild :menu="menu" />
            </VuePerfectScrollbar>
        </div>

    </div>
</template>

<script>
    import {SidebarMenu} from 'vue-sidebar-menu'
    import VuePerfectScrollbar from 'vue-perfect-scrollbar'
    import ApiService from '../../Services/apiService'

    export default {
        components: {
            SidebarMenu,
            VuePerfectScrollbar
        },
        computed: {

        },
        async beforeMount() {
            this.menu = await this.getMenu();
        },
        data() {
            return {
                isOpen: false,
                sidebarActive: false,
                apiService: new ApiService(),
                menu: null,
                collapsed: true,
                windowWidth: 0,
            }
        },
        props: {
            sidebarbg: String,
        },
        methods: {
            async getMenu() {
                const namespaceList = await this.apiService.listNamespace();
                const children = namespaceList.map(item => {
                    return  {
                        href: `/namespaces/${item}`,
                        title: item
                    }
                })
                return [
                    {
                        header: true,
                        title: 'Navigation',
                    },
                    {
                        title: 'Namespaces',
                        icon: 'pe-7s-rocket',
                        child: children
                    },
                ];
            },
            toggleBodyClass(className) {
                const el = document.body;
                this.isOpen = !this.isOpen;

                if (this.isOpen) {
                    el.classList.add(className);
                } else {
                    el.classList.remove(className);
                }
            },
            toggleSidebarHover(add, className) {
                const el = document.body;
                this.sidebarActive = !this.sidebarActive;

                this.windowWidth = document.documentElement.clientWidth;

                if (this.windowWidth > '992') {
                    if (add === 'add') {
                        el.classList.add(className);
                    } else {
                        el.classList.remove(className);
                    }
                }
            },
            getWindowWidth() {
                const el = document.body;

                this.windowWidth = document.documentElement.clientWidth;

                if (this.windowWidth < '1350') {
                    el.classList.add('closed-sidebar', 'closed-sidebar-md');
                } else {
                    el.classList.remove('closed-sidebar', 'closed-sidebar-md');
                }
            },
        },
        mounted() {
            this.$nextTick(function () {
                window.addEventListener('resize', this.getWindowWidth);

                //Init
                this.getWindowWidth()
            })
        },

        beforeDestroy() {
            window.removeEventListener('resize', this.getWindowWidth);
        }
    }
</script>
