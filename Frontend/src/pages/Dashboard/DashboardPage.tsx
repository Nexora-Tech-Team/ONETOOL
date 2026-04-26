import { useEffect, useState } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { RootState, AppDispatch } from '@/store'
import { fetchMe } from '@/store/slices/authSlice'
import { dashboardService, taskService, projectService, teamService } from '@/services/api'
import { Link } from 'react-router-dom'
import { Clock, CheckSquare, Calendar, CreditCard, FolderKanban, Users } from 'lucide-react'
import { PieChart, Pie, Cell, Tooltip } from 'recharts'
import { Loading, StatusBadge, ProgressBar } from '@/components/common'
import { toast } from 'react-toastify'
import { formatIDR } from '@/utils/format'
import clsx from 'clsx'

const TASK_COLORS = ['#f97316', '#3b82f6', '#10b981', '#ef4444']

export default function DashboardPage() {
  const { user } = useSelector((s: RootState) => s.auth)
  const dispatch = useDispatch<AppDispatch>()
  const [stats, setStats] = useState<any>(null)
  const [tasks, setTasks] = useState<any[]>([])
  const [projects, setProjects] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [clockLoading, setClockLoading] = useState(false)

  const loadData = (uid?: number) => {
    Promise.all([
      dashboardService.getStats(),
      taskService.list({ assigned_to_id: uid, limit: 5 }),
      projectService.list({ status: 'open', limit: 3 }),
    ]).then(([s, t, p]) => {
      setStats(s.data)
      setTasks(t.data.data || [])
      setProjects(p.data.data || [])
    }).finally(() => setLoading(false))
  }

  useEffect(() => { loadData(user?.id) }, [user?.id])

  const handleClock = async () => {
    setClockLoading(true)
    try {
      if (user?.clocked_in) {
        await teamService.clockOut()
        toast.success('Clocked out successfully!')
      } else {
        await teamService.clockIn()
        toast.success('Clocked in successfully!')
      }
      await dispatch(fetchMe())
      const s = await dashboardService.getStats()
      setStats(s.data)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Clock action failed')
    } finally {
      setClockLoading(false)
    }
  }

  if (loading) return <div className="p-6"><Loading /></div>

  const taskChartData = [
    { name: 'To do', value: stats?.tasks_todo ?? 0 },
    { name: 'In progress', value: stats?.tasks_in_progress ?? 0 },
    { name: 'Done', value: stats?.tasks_done ?? 0 },
    { name: 'Expired', value: stats?.tasks_expired ?? 0 },
  ]

  const incomeExpenseData = [
    { value: stats?.total_income ?? 0 },
    { value: stats?.total_expenses ?? 0 },
  ]
  const totalIE = (stats?.total_income ?? 0) + (stats?.total_expenses ?? 0)
  const incomePercent = totalIE > 0 ? Math.round((stats?.total_income / totalIE) * 100) : 0

  const invoiceRows = [
    { label: 'Overdue',        color: '#ef4444', val: stats?.overdue_amount ?? 0 },
    { label: 'Not paid',       color: '#f97316', val: stats?.not_paid_amount ?? 0 },
    { label: 'Partially paid', color: '#3b82f6', val: stats?.partially_paid_amount ?? 0 },
    { label: 'Fully paid',     color: '#1d4ed8', val: stats?.fully_paid_amount ?? 0 },
    { label: 'Draft',          color: '#9ca3af', val: stats?.draft_amount ?? 0 },
  ]

  const clockedOut = (stats?.total_members ?? 0) - (stats?.clocked_in_count ?? 0)

  return (
    <div className="p-5 space-y-4">
      {/* Stats Row */}
      <div className="grid grid-cols-4 gap-3">
        {/* Clock In/Out */}
        <div
          className="rounded-lg p-4 flex flex-col gap-1.5 transition-all duration-500"
          style={{
            background: user?.clocked_in
              ? 'linear-gradient(135deg,#10b981,#059669)'
              : 'linear-gradient(135deg,#ec4899,#f43f5e)',
            color: 'white',
          }}
        >
          <div className="flex items-center gap-2 text-xs opacity-90">
            <Clock size={14} /> Clock In/Out
          </div>
          <div className="text-xs opacity-80">
            {user?.clocked_in ? 'You are currently clocked in' : 'You are currently clocked out'}
          </div>
          <button
            onClick={handleClock}
            disabled={clockLoading}
            className="mt-1 text-xs bg-white/20 hover:bg-white/30 px-3 py-1 rounded-full w-fit transition-colors disabled:opacity-60"
          >
            {clockLoading ? '...' : user?.clocked_in ? '⏹ Clock Out' : '+ Clock In'}
          </button>
        </div>
        <div className="stat-card">
          <div className="stat-icon bg-blue-50"><CheckSquare size={20} className="text-blue-600" /></div>
          <div>
            <div className="stat-val">{stats?.open_tasks ?? 0}</div>
            <div className="stat-label">My open tasks</div>
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-icon bg-indigo-50"><Calendar size={20} className="text-indigo-600" /></div>
          <div>
            <div className="stat-val">{stats?.on_leave_today ?? 0}</div>
            <div className="stat-label">On leave today</div>
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-icon bg-rose-50"><CreditCard size={20} className="text-rose-500" /></div>
          <div>
            <div className="stat-val text-sm">{formatIDR(stats?.due_amount ?? 0)}</div>
            <div className="stat-label">Due</div>
          </div>
        </div>
      </div>

      {/* Row 2 */}
      <div className="grid grid-cols-3 gap-3">
        {/* Projects Overview */}
        <div className="card">
          <div className="card-header">
            <span className="section-title flex items-center gap-1.5">
              <FolderKanban size={13} className="text-gray-400" /> Projects Overview
            </span>
          </div>
          <div className="card-body">
            <div className="flex gap-5 mb-3">
              <div className="text-center"><div className="text-lg font-semibold text-blue-600">{stats?.open_projects ?? 0}</div><div className="text-xs text-gray-400">Open</div></div>
              <div className="text-center"><div className="text-lg font-semibold text-green-600">{stats?.completed_projects ?? 0}</div><div className="text-xs text-gray-400">Completed</div></div>
              <div className="text-center"><div className="text-lg font-semibold text-orange-500">{stats?.hold_projects ?? 0}</div><div className="text-xs text-gray-400">Hold</div></div>
            </div>
            {(() => {
              const total = (stats?.open_projects ?? 0) + (stats?.completed_projects ?? 0) + (stats?.hold_projects ?? 0)
              const prog = total > 0 ? Math.round(((stats?.completed_projects ?? 0) / total) * 100) : 0
              return <>
                <div className="text-xs text-gray-400 mb-1">Completion {prog}%</div>
                <ProgressBar value={prog} />
              </>
            })()}
          </div>
        </div>

        {/* Invoice Overview */}
        <div className="card">
          <div className="card-header"><span className="section-title">Invoice Overview</span></div>
          <div className="card-body space-y-1.5">
            {invoiceRows.map(item => (
              <div key={item.label} className="flex items-center justify-between text-xs">
                <div className="flex items-center gap-1.5">
                  <div className="w-2 h-2 rounded-full" style={{ background: item.color }} />
                  <span className="text-gray-500">{item.label}</span>
                </div>
                <span className="font-medium">{formatIDR(item.val)}</span>
              </div>
            ))}
            <div className="pt-2 border-t border-gray-100 text-xs flex justify-between">
              <span className="text-gray-400">Total invoiced</span>
              <span className="font-semibold">{formatIDR(stats?.total_invoiced ?? 0)}</span>
            </div>
          </div>
        </div>

        {/* Income vs Expenses */}
        <div className="card">
          <div className="card-header"><span className="section-title">Income vs Expenses</span></div>
          <div className="card-body flex items-center gap-3">
            <PieChart width={80} height={80}>
              <Pie data={incomeExpenseData} dataKey="value" cx={35} cy={35} innerRadius={22} outerRadius={35} strokeWidth={2}>
                <Cell fill="#1d4ed8" /><Cell fill="#ef4444" />
              </Pie>
              <Tooltip formatter={(v: any) => formatIDR(v)} contentStyle={{ fontSize: 10 }} />
            </PieChart>
            <div className="text-xs space-y-1">
              <div className="font-medium text-gray-400 mb-1">Total ({incomePercent}% income)</div>
              <div className="flex items-center gap-1.5"><div className="w-2 h-2 rounded-full bg-blue-700" /><span>{formatIDR(stats?.total_income ?? 0)}</span></div>
              <div className="flex items-center gap-1.5"><div className="w-2 h-2 rounded-full bg-red-500" /><span>{formatIDR(stats?.total_expenses ?? 0)}</span></div>
            </div>
          </div>
        </div>
      </div>

      {/* Row 3 */}
      <div className="grid grid-cols-3 gap-3">
        {/* Tasks Overview */}
        <div className="card">
          <div className="card-header"><span className="section-title">All Tasks Overview</span></div>
          <div className="card-body flex items-center gap-4">
            <PieChart width={90} height={90}>
              <Pie data={taskChartData} dataKey="value" cx={40} cy={40} innerRadius={24} outerRadius={40} strokeWidth={2}>
                {taskChartData.map((_, i) => <Cell key={i} fill={TASK_COLORS[i]} />)}
              </Pie>
            </PieChart>
            <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs">
              {taskChartData.map((item, i) => (
                <div key={item.name} className="flex items-center gap-1.5">
                  <div className="w-2 h-2 rounded-full" style={{ background: TASK_COLORS[i] }} />
                  <span className="text-gray-500">{item.name}</span>
                  <span className="font-semibold ml-auto">{item.value}</span>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Team Overview */}
        <div className="card">
          <div className="card-header"><span className="section-title flex items-center gap-1.5"><Users size={13} className="text-gray-400" />Team Members Overview</span></div>
          <div className="card-body">
            <div className="grid grid-cols-2 gap-2 mb-3">
              <div className="bg-gray-50 rounded-lg p-2 text-center">
                <div className="text-xl font-semibold">{stats?.total_members ?? 0}</div>
                <div className="text-xs text-gray-400">Team members</div>
              </div>
              <div className="bg-gray-50 rounded-lg p-2 text-center">
                <div className="text-xl font-semibold text-orange-500">{stats?.on_leave_today ?? 0}</div>
                <div className="text-xs text-gray-400">On leave today</div>
              </div>
            </div>
            <div className="flex gap-2">
              <div className="flex-1 bg-red-50 rounded-lg p-2 text-center">
                <div className="text-base font-semibold text-red-500">{stats?.clocked_in_count ?? 0}</div>
                <div className="text-xs text-red-400">Clocked In</div>
              </div>
              <div className="flex-1 bg-blue-50 rounded-lg p-2 text-center">
                <div className="text-base font-semibold text-blue-600">{clockedOut}</div>
                <div className="text-xs text-blue-400">Clocked Out</div>
              </div>
            </div>
          </div>
        </div>

        {/* Timesheet placeholder */}
        <div className="card">
          <div className="card-header">
            <span className="section-title">My Timesheet</span>
            <Link to="/team/timecards" className="text-xs text-blue-600 hover:underline">View all</Link>
          </div>
          <div className="card-body flex items-center justify-center h-20">
            <span className="text-xs text-gray-400">View detailed timesheet in Team → Time Cards</span>
          </div>
        </div>
      </div>

      {/* Row 4 */}
      <div className="grid grid-cols-3 gap-3">
        {/* My Tasks */}
        <div className="card col-span-2">
          <div className="card-header">
            <span className="section-title">My Tasks</span>
            <Link to="/tasks" className="text-xs text-blue-600 hover:underline">View all</Link>
          </div>
          <div>
            <table className="table">
              <thead><tr><th>ID</th><th>Title</th><th>Start date</th><th>Deadline</th><th>Status</th></tr></thead>
              <tbody>
                {tasks.length === 0 ? (
                  <tr><td colSpan={5} className="text-center text-gray-400 py-6 text-xs">No tasks found</td></tr>
                ) : tasks.map(task => (
                  <tr key={task.id}>
                    <td className="text-gray-400">{task.id}</td>
                    <td><Link to={`/tasks`} className="text-blue-600 hover:underline">{task.title}</Link></td>
                    <td className="text-gray-400">{task.start_date ? new Date(task.start_date).toLocaleDateString('id') : '-'}</td>
                    <td className={clsx('text-gray-400', new Date(task.deadline) < new Date() && 'text-red-500')}>
                      {task.deadline ? new Date(task.deadline).toLocaleDateString('id') : '-'}
                    </td>
                    <td><StatusBadge status={task.status} /></td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Open Projects */}
        <div className="card">
          <div className="card-header">
            <span className="section-title">Open Projects</span>
            <Link to="/projects" className="text-xs text-blue-600 hover:underline">View all</Link>
          </div>
          <div className="card-body space-y-3">
            {projects.length === 0 ? (
              <div className="text-xs text-gray-400 text-center py-4">No open projects</div>
            ) : projects.map(p => (
              <div key={p.id}>
                <div className="flex items-center justify-between mb-1">
                  <Link to={`/projects/${p.id}`} className="text-xs text-blue-600 hover:underline truncate flex-1 mr-2">{p.title}</Link>
                  <span className="text-xs text-gray-400">{p.progress}%</span>
                </div>
                <ProgressBar value={p.progress} />
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
