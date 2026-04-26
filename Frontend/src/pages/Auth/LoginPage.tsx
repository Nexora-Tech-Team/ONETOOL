import { useState, useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import { login } from '@/store/slices/authSlice'
import { RootState, AppDispatch } from '@/store'
import { Eye, EyeOff } from 'lucide-react'
import logoUrl from '../../../logo/logo_ref.png'
import loginBg from '../../../logo/login-bg.png'
import nexoraUrl from '../../../logo/Logo_Nexora_Part.png'

export default function LoginPage() {
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()
  const { token, loading, error } = useSelector((s: RootState) => s.auth)

  const [email,    setEmail]    = useState(() => localStorage.getItem('remembered_email') || '')
  const [password, setPassword] = useState('')
  const [showPass, setShowPass] = useState(false)
  const [remember, setRemember] = useState(() => !!localStorage.getItem('remembered_email'))

  useEffect(() => { if (token) navigate('/dashboard') }, [token])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (remember) localStorage.setItem('remembered_email', email)
    else          localStorage.removeItem('remembered_email')
    dispatch(login({ email, password, remember }))
  }

  return (
    <div className="min-h-screen flex bg-white">

      {/* ── Left panel ──────────────────────────────────────── */}
      <div
        className="w-full lg:w-1/2 xl:w-5/12 2xl:w-1/3 flex flex-col justify-center relative"
        style={{ borderTop: '3px solid #2aacb8' }}
      >
        <div className="w-full px-8 py-10">

          {/* Logo */}
          <div className="flex justify-center mb-6">
            <img src={logoUrl} alt="CBQA Global" style={{ width: 160 }} />
          </div>

          {/* Title */}
          <div className="text-center mb-5">
            <h2 className="text-[1.6rem] font-bold text-gray-900 leading-tight">
              Welcome back
            </h2>
            <p className="text-sm text-gray-500 mt-1">
              Sign in to access your{' '}
              <span className="text-[#2aacb8] font-medium">IT Audit System</span>
            </p>
          </div>

          {/* Error */}
          {error && (
            <div className="text-sm text-red-600 bg-red-50 border border-red-200 rounded px-3 py-2 mb-4">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} autoComplete="off">

            {/* Email */}
            <div className="mb-3">
              <label className="block text-sm text-gray-700 mb-1">
                Email or Username
              </label>
              <input
                type="text"
                value={email}
                onChange={e => setEmail(e.target.value)}
                placeholder="Enter your email or username"
                required
                autoComplete="off"
                style={{ backgroundColor: '#f0f5fb' }}
                className="w-full px-3 py-2.5 text-sm border border-[#c8d3de] rounded
                           text-gray-900 placeholder-gray-400
                           focus:outline-none focus:border-[#2aacb8] focus:ring-1 focus:ring-[#2aacb8]
                           transition-colors"
              />
            </div>

            {/* Password */}
            <div className="mb-2">
              <div className="flex items-center justify-between mb-1">
                <label className="text-sm text-gray-700">Password</label>
                <button
                  type="button"
                  onClick={() => navigate('/forgot-password')}
                  className="text-sm text-[#2aacb8] hover:text-[#1d8a94] transition-colors"
                >
                  I forgot password
                </button>
              </div>
              <div
                className="flex rounded border border-[#c8d3de] overflow-hidden
                            focus-within:border-[#2aacb8] focus-within:ring-1 focus-within:ring-[#2aacb8]
                            transition-colors"
                style={{ backgroundColor: '#f0f5fb' }}
              >
                <input
                  type={showPass ? 'text' : 'password'}
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  placeholder="Your password"
                  required
                  autoComplete="off"
                  className="flex-1 px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400
                             focus:outline-none border-0 bg-transparent"
                />
                <button
                  type="button"
                  onClick={() => setShowPass(!showPass)}
                  className="px-3 flex items-center text-gray-400 hover:text-gray-600 bg-transparent border-0"
                >
                  {showPass ? <EyeOff size={15} /> : <Eye size={15} />}
                </button>
              </div>
            </div>

            {/* Remember me */}
            <div className="flex items-center gap-2 mb-5">
              <input
                type="checkbox"
                id="remember"
                checked={remember}
                onChange={e => setRemember(e.target.checked)}
                className="w-4 h-4 rounded border-gray-300 accent-[#2aacb8] cursor-pointer"
              />
              <label htmlFor="remember" className="text-sm text-gray-600 cursor-pointer select-none">
                Remember me on this device
              </label>
            </div>

            {/* Submit */}
            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 rounded text-sm font-medium text-white
                         bg-[#2aacb8] hover:bg-[#239dab] active:bg-[#1d8a94]
                         flex items-center justify-center transition-colors
                         disabled:opacity-70"
            >
              {loading
                ? <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                : 'Sign in'}
            </button>
          </form>

        </div>

        {/* Nexora branding – bottom-left */}
        <div className="absolute bottom-0 left-4">
          <img src={nexoraUrl} alt="Nexora – Part of CBQA Global Group" style={{ width: 130, opacity: 0.90 }} />
        </div>

      </div>

      {/* ── Right panel ─────────────────────────────────────── */}
      <div className="hidden lg:block flex-1 relative overflow-hidden">
        {/* Blurred background layer fills gap areas */}
        <div
          className="absolute inset-0 bg-center bg-cover"
          style={{
            backgroundImage: `url(${loginBg})`,
            filter: 'blur(24px)',
            transform: 'scale(1.08)',
          }}
        />
        {/* Main image layer, contained without cropping */}
        <div
          className="absolute inset-0 bg-center bg-no-repeat"
          style={{
            backgroundImage: `url(${loginBg})`,
            backgroundSize: 'contain',
          }}
        />
      </div>

    </div>
  )
}
