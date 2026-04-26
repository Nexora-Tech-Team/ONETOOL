import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { useState } from 'react'
import { RootState } from '@/store'
import { logout, canRead } from '@/store/slices/authSlice'
import { toggleSidebar } from '@/store/slices/uiSlice'
import logoUrl from '../../../logo/logo_ref.png'
import {
  LayoutDashboard, Calendar, Users, FolderKanban, CheckSquare,
  TrendingUp, CreditCard, FileText, MessageSquare, UserCircle,
  FolderOpen, Receipt, BarChart2, CheckCircle, Menu, Search,
  Bell, Globe, Clock, Plus, LogOut, ChevronDown, ChevronRight, Settings
} from 'lucide-react'
import clsx from 'clsx'

const navItems = [
  { to: '/dashboard', icon: LayoutDashboard, label: 'Dashboard', menu: 'dashboard' },
  { to: '/events', icon: Calendar, label: 'Events', menu: 'events' },
  { to: '/clients', icon: Users, label: 'Clients', menu: 'clients' },
  { to: '/projects', icon: FolderKanban, label: 'Projects', menu: 'projects' },
  { to: '/tasks', icon: CheckSquare, label: 'Tasks', menu: 'tasks' },
  { to: '/leads', icon: TrendingUp, label: 'Leads', menu: 'leads' },
  {
    label: 'Sales', icon: CreditCard, menu: 'sales', children: [
      { to: '/sales/invoices', label: 'Invoices' },
      { to: '/sales/orders', label: 'Orders list' },
      { to: '/sales/store', label: 'Store' },
      { to: '/sales/payments', label: 'Payments' },
      { to: '/sales/items', label: 'Items' },
      { to: '/sales/contracts', label: 'Contracts' },
    ]
  },
  { to: '/notes', icon: FileText, label: 'Notes', menu: 'notes' },
  { to: '/messages', icon: MessageSquare, label: 'Messages', menu: 'messages' },
  {
    label: 'Team', icon: UserCircle, menu: 'team', children: [
      { to: '/team/members', label: 'Team members' },
      { to: '/team/timecards', label: 'Time cards' },
      { to: '/team/leave', label: 'Leave' },
      { to: '/team/announcements', label: 'Announcements' },
      { to: '/team/help', label: 'Help' },
    ]
  },
  { to: '/files', icon: FolderOpen, label: 'Files', menu: 'files' },
  { to: '/expenses', icon: Receipt, label: 'Expenses', menu: 'expenses' },
  { to: '/reports', icon: BarChart2, label: 'Reports', menu: 'reports' },
  { to: '/todo', icon: CheckCircle, label: 'To do', menu: 'todo' },
  {
    label: 'Settings', icon: Settings, menu: 'settings', children: [
      { to: '/settings/users', label: 'User Accounts' },
      { to: '/settings/roles', label: 'Roles' },
      { to: '/settings/audit-log', label: 'Audit Trail' },
    ]
  },
]

export default function Layout() {
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const { user, permissions } = useSelector((s: RootState) => s.auth)
  const { sidebarOpen } = useSelector((s: RootState) => s.ui)
  const [expanded, setExpanded] = useState<string[]>(['Sales', 'Team'])

  const visibleNav = navItems.filter(item => canRead(permissions, user?.role, item.menu))

  const toggleExpand = (label: string) => {
    setExpanded(prev =>
      prev.includes(label) ? prev.filter(l => l !== label) : [...prev, label]
    )
  }

  const handleLogout = () => {
    dispatch(logout())
    navigate('/login')
  }

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Sidebar */}
      <aside className={clsx(
        'flex flex-col flex-shrink-0 transition-all duration-200 overflow-y-auto',
        sidebarOpen ? 'w-[190px]' : 'w-0 overflow-hidden'
      )} style={{ backgroundColor: '#073d47' }}>
        {/* Logo */}
        <div className="px-4 py-4 flex items-center justify-center" style={{ borderBottom: '1px solid rgba(255,255,255,0.1)' }}>
          <img src={logoUrl} alt="CBQA Global" style={{ width: 130, filter: 'brightness(0) invert(1)' }} />
        </div>

        {/* Nav */}
        <nav className="flex-1 py-2">
          {visibleNav.map((item) => {
            if ('children' in item) {
              const isExp = expanded.includes(item.label)
              const Icon = item.icon
              return (
                <div key={item.label}>
                  <button
                    onClick={() => toggleExpand(item.label)}
                    className="nav-link w-full text-left"
                  >
                    <Icon size={15} className="flex-shrink-0" />
                    <span className="flex-1">{item.label}</span>
                    {isExp ? <ChevronDown size={12} /> : <ChevronRight size={12} />}
                  </button>
                  {isExp && item.children?.map(child => (
                    <NavLink
                      key={child.to}
                      to={child.to}
                      className={({ isActive }) => clsx('nav-link nav-sub', isActive && 'active')}
                    >
                      {child.label}
                    </NavLink>
                  ))}
                </div>
              )
            }
            const Icon = item.icon
            return (
              <NavLink
                key={item.to}
                to={item.to}
                className={({ isActive }) => clsx('nav-link', isActive && 'active')}
              >
                <Icon size={15} className="flex-shrink-0" />
                {item.label}
              </NavLink>
            )
          })}
        </nav>

        {/* User at bottom */}
        <div className="p-3" style={{ borderTop: '1px solid rgba(255,255,255,0.1)' }}>
          <div className="px-2 pb-2 text-[10px]" style={{ color: 'rgba(255,255,255,0.3)' }}>
            v1.0.2
          </div>
          <button
            onClick={handleLogout}
            className="flex items-center gap-2 w-full px-2 py-1.5 text-xs rounded transition-colors"
            style={{ color: 'rgba(255,255,255,0.55)' }}
            onMouseEnter={e => (e.currentTarget.style.color = '#fff')}
            onMouseLeave={e => (e.currentTarget.style.color = 'rgba(255,255,255,0.55)')}
          >
            <LogOut size={13} />
            Sign out
          </button>
        </div>
      </aside>

      {/* Main */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Topbar */}
        <header className="bg-white border-b border-gray-200 h-12 flex items-center px-4 gap-3 flex-shrink-0">
          <button
            onClick={() => dispatch(toggleSidebar())}
            className="w-8 h-8 flex items-center justify-center rounded-md text-gray-400 hover:bg-gray-100"
          >
            <Menu size={16} />
          </button>

          <div className="flex items-center gap-2 ml-auto">
            {/* Icons */}
            {[Search, Plus, Globe, Clock, Bell].map((Icon, i) => (
              <button key={i} className="w-8 h-8 flex items-center justify-center rounded-md text-gray-400 hover:bg-gray-100 relative">
                <Icon size={15} />
                {i === 4 && (
                  <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full border-2 border-white" />
                )}
              </button>
            ))}

            {/* User chip */}
            <div className="flex items-center gap-2 px-2.5 py-1.5 bg-gray-50 rounded-full cursor-pointer hover:bg-gray-100 ml-1">
              <div className="w-6 h-6 rounded-full bg-blue-800 text-white flex items-center justify-center text-[10px] font-semibold">
                {user?.name?.split(' ').map(w => w[0]).join('').slice(0, 2).toUpperCase()}
              </div>
              <span className="text-xs font-medium text-gray-700">{user?.name}</span>
            </div>
          </div>
        </header>

        {/* Page content */}
        <main className="flex-1 overflow-y-auto">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
