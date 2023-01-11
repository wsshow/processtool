<template>
    <div class="about">
        <el-card shadow="always">
            <el-row style="margin-bottom: 10px">
                <el-col :span="16">
                    <el-input
                        v-model="inputVal"
                        placeholder="请输入搜索内容"
                        class="input-with-select"
                    >
                        <template #prepend>
                            <el-select
                                v-model="selectOp"
                                placeholder="Select"
                                style="width: 90px"
                            >
                                <el-option
                                    v-for="item in options"
                                    :key="item.value"
                                    :label="item.label"
                                    :value="item.value"
                                >
                                </el-option>
                            </el-select>
                        </template>
                        <template #append>
                            <el-button :icon="Search" />
                        </template>
                    </el-input>
                </el-col>
                <el-col :span="8">
                    <el-tag size="large">服务器响应时间：{{ requestDuration }}</el-tag>
                    <el-tag size="large" type="success">进程数：{{ procsCount }}</el-tag>
                    <el-tag size="large" type="info">数据更新时间：{{ dataRefreshTime }}</el-tag>
                </el-col>
            </el-row>
            <el-table
                v-loading="loading"
                :data="filterTableData"
                :default-sort="{ prop: 'cpuPercent', order: 'descending' }"
                :height="tableHeight"
                :border="true"
                style="width: 100%"
            >
                <el-table-column type="expand">
                    <template #default="props">
                        <el-card>
                            <el-table
                                :data="props.row.children"
                                :default-sort="{
                                    prop: 'status',
                                    order: 'descending',
                                }"
                            >
                                <el-table-column
                                    v-for="{ prop, label } in childrenCols"
                                    :key="prop"
                                    :prop="prop"
                                    :label="label"
                                    :show-overflow-tooltip="true"
                                    resizable
                                >
                                </el-table-column>
                            </el-table>
                        </el-card>
                    </template>
                </el-table-column>
                <el-table-column
                    v-for="{ prop, label } in cols"
                    :key="prop"
                    :prop="prop"
                    :label="label"
                    :show-overflow-tooltip="true"
                    :sortable="prop !== 'cmdline'"
                >
                </el-table-column>
            </el-table>
        </el-card>
    </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref, computed } from 'vue'
import { Search } from '@element-plus/icons-vue'
import axios from 'axios'
interface rowChildrenElement {
  localaddr: string;
  remoteaddr: string;
  protocol: string;
  status: string;
}
interface rowElement {
  pid: number;
  name: string;
  cmdline: string;
  cpuPercent: string;
  memPercent: string;
  memSize: string;
  children?: rowChildrenElement[];
}
class Timer {
  handler: any
  timeout?: number
  id: number | undefined
  constructor (handler: any, timeout?: number) {
      this.handler = handler
      this.timeout = timeout
  }

  Start = () => {
      const interval = () => {
          this.handler()
          this.id = setTimeout(interval, this.timeout)
      }

      this.id = setTimeout(interval, this.timeout)
  }

   Stop = () => {
       clearTimeout(this.id)
   }
}
export default defineComponent({
    name: 'HomeView',
    setup () {
        const loading = ref(true)
        const cols = [
            { prop: 'pid', label: '进程ID' },
            { prop: 'name', label: '进程名' },
            { prop: 'cmdline', label: '命令行参数' },
            { prop: 'cpuPercent', label: 'CPU占用率' },
            { prop: 'memPercent', label: '内存占用率' },
            { prop: 'memSize', label: '内存占用大小' }
        ]
        const childrenCols = [
            { prop: 'localaddr', label: '本地地址' },
            { prop: 'remoteaddr', label: '远程地址' },
            { prop: 'protocol', label: '连接协议' },
            { prop: 'status', label: '连接状态' }
        ]
        const options = [
            {
                value: 'name',
                label: '进程名'
            },
            {
                value: 'pid',
                label: '进程ID'
            }
        ]

        const http = axios.create({
            baseURL: 'http://127.0.0.1:9091',
            timeout: 10000
        })

        http.interceptors.request.use(
            function (config: any) {
                config.metadata = { startTime: new Date() }
                return config
            },
            function (error: any) {
                return Promise.reject(error)
            }
        )

        http.interceptors.response.use(
            function (response: any) {
                response.config.metadata.endTime = new Date()
                response.duration =
          response.config.metadata.endTime -
          response.config.metadata.startTime
                return response
            },
            function (error: any) {
                error.config.metadata.endTime = new Date()
                error.duration =
          error.config.metadata.endTime -
          error.config.metadata.startTime
                return Promise.reject(error)
            }
        )

        function getProtocol (family: number, type: number): string {
            if (family === 2) {
                if (type === 1) {
                    return 'tcp'
                }
                if (type === 2) {
                    return 'udp'
                }
            }
            if (family === 23) {
                if (type === 1) {
                    return 'tcp6'
                }
                if (type === 2) {
                    return 'udp6'
                }
            }
            return 'unknown'
        }

        function getChildren (e: any): rowChildrenElement[] {
            const conns = e.connections
            const len = conns.length
            const rows: rowChildrenElement[] = []
            for (let index = 0; index < len; index++) {
                const conn = conns[index]
                rows.push({
                    localaddr: conn.localaddr.ip + ':' + conn.localaddr.port,
                    remoteaddr: conn.remoteaddr.ip + ':' + conn.remoteaddr.port,
                    protocol: getProtocol(conn.family, conn.type),
                    status: conn.status
                })
            }
            return rows
        }

        const tableData: rowElement[] = reactive([])
        const requestDuration = ref('')
        const procsCount = ref()
        const dataRefreshTime = ref('')

        function refreshData () {
            http.get('/procs')
                .then(function (response: any) {
                    const data = response.data.data
                    if (data.length > 0) {
                        tableData.length = 0
                    }
                    for (const e of data) {
                        if (!('connections' in e)) {
                            continue
                        }
                        tableData.push({
                            pid: e.pid,
                            name: e.name,
                            cmdline: e.cmd_line,
                            cpuPercent: e.cpu_percent.toFixed(2) + '%',
                            memPercent: e.mem_percent.toFixed(2) + '%',
                            memSize: (e.memory_info.vms >> 20).toFixed(2) + 'MB',
                            children: getChildren(e)
                        })
                    }
                    loading.value = false
                    requestDuration.value = response.duration / 1000 + 's'
                    procsCount.value = tableData.length
                    dataRefreshTime.value = new Date().toLocaleString()
                })
                .catch(function (error: any) {
                    console.log(error)
                })
        }

        const timer = new Timer(refreshData, 5000)
        timer.Start()

        setTimeout(timer.Stop, 10000)

        const inputVal = ref('')
        const selectOp = ref(options[0].value)
        const filterTableData = computed(() =>
            tableData.filter((data) => {
                if (!inputVal.value) {
                    return true
                }
                switch (selectOp.value) {
                case 'name':
                    return data.name
                        .toLowerCase()
                        .includes(inputVal.value.toLowerCase())
                case 'pid':
                    return data.pid.toString().includes(inputVal.value)
                default:
                    return data.name
                        .toLowerCase()
                        .includes(inputVal.value.toLowerCase())
                }
            })
        )

        const tableHeight = ref(window.screen.availHeight * 0.82)

        return {
            cols,
            childrenCols,
            filterTableData,
            loading,
            requestDuration,
            procsCount,
            dataRefreshTime,
            inputVal,
            Search,
            selectOp,
            options,
            tableHeight
        }
    }
})
</script>
