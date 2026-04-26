import { useNavigate } from 'react-router-dom'
import logoUrl from '../../../logo/logo_ref.png'
import loginBg from '../../../logo/login-bg.png'
import nexoraUrl from '../../../logo/Logo_Nexora_Part.png'
import { ArrowLeft, ShieldAlert } from 'lucide-react'

export default function ForgotPasswordPage() {
  const navigate = useNavigate()

  return (
    <div className="min-h-screen flex bg-white">

      {/* ── Left panel ──────────────────────────────────────── */}
      <div
        className="w-full lg:w-1/2 xl:w-5/12 2xl:w-1/3 flex flex-col justify-center relative"
        style={{ borderTop: '3px solid #2aacb8' }}
      >
        <div className="w-full px-8 py-10">

          {/* Logo */}
          <div className="flex justify-center mb-8">
            <img src={logoUrl} alt="CBQA Global" style={{ width: 160 }} />
          </div>

          {/* Icon */}
          <div className="flex justify-center mb-5">
            <div className="w-16 h-16 rounded-full flex items-center justify-center"
              style={{ backgroundColor: '#e6f7f8' }}>
              <ShieldAlert size={30} style={{ color: '#2aacb8' }} />
            </div>
          </div>

          {/* Message */}
          <h2 className="text-center text-[1.4rem] font-bold text-gray-900 mb-2">
            Forgot your password?
          </h2>
          <p className="text-center text-sm text-gray-500 mb-1 leading-relaxed">
            For security reasons, password resets are managed by your system administrator.
          </p>
          <p className="text-center text-sm text-gray-500 mb-6 leading-relaxed">
            Please contact your <span className="font-medium text-[#2aacb8]">admin</span> to reset your password.
          </p>

          {/* Back button */}
          <button
            onClick={() => navigate('/login')}
            className="w-full py-2.5 rounded text-sm font-medium text-white flex items-center justify-center gap-2 transition-colors"
            style={{ backgroundColor: '#2aacb8' }}
            onMouseEnter={e => (e.currentTarget.style.backgroundColor = '#239dab')}
            onMouseLeave={e => (e.currentTarget.style.backgroundColor = '#2aacb8')}
          >
            <ArrowLeft size={15} />
            Back to Login
          </button>
        </div>

        {/* Nexora branding – bottom-left */}
        <div className="absolute bottom-0 left-4">
          <img src={nexoraUrl} alt="Nexora – Part of CBQA Global Group" style={{ width: 130, opacity: 0.90 }} />
        </div>
      </div>

      {/* ── Right panel ─────────────────────────────────────── */}
      <div className="hidden lg:block flex-1 relative overflow-hidden">
        <div
          className="absolute inset-0 bg-center bg-cover"
          style={{
            backgroundImage: `url(${loginBg})`,
            filter: 'blur(24px)',
            transform: 'scale(1.08)',
          }}
        />
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
