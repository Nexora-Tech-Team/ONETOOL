import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit'
import { authService } from '@/services/api'

interface User {
  id: number
  name: string
  email: string
  job_title: string
  role: string
  avatar: string
  clocked_in: boolean
}

interface AuthState {
  user: User | null
  token: string | null
  loading: boolean
  error: string | null
}

const storedToken = localStorage.getItem('token') || sessionStorage.getItem('token')
const storedUser = localStorage.getItem('user') || sessionStorage.getItem('user')

const initialState: AuthState = {
  user: JSON.parse(storedUser || 'null'),
  token: storedToken,
  loading: false,
  error: null,
}

export const login = createAsyncThunk(
  'auth/login',
  async (
    { email, password, remember }: { email: string; password: string; remember: boolean },
    { rejectWithValue }
  ) => {
    try {
      const res = await authService.login(email, password)
      return { ...res.data, remember }
    } catch (err: any) {
      return rejectWithValue(err.response?.data?.error || 'Login failed')
    }
  }
)

export const fetchMe = createAsyncThunk('auth/me', async (_, { rejectWithValue }) => {
  try {
    const res = await authService.me()
    return res.data
  } catch (err: any) {
    return rejectWithValue(err.response?.data?.error || 'Failed to fetch user')
  }
})

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    logout(state) {
      state.user = null
      state.token = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      sessionStorage.removeItem('token')
      sessionStorage.removeItem('user')
    },
    setUser(state, action: PayloadAction<User>) {
      state.user = action.payload
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(login.pending, (state) => { state.loading = true; state.error = null })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false
        state.token = action.payload.token
        state.user = action.payload.user
        if (action.payload.remember) {
          localStorage.setItem('token', action.payload.token)
          localStorage.setItem('user', JSON.stringify(action.payload.user))
          sessionStorage.removeItem('token')
          sessionStorage.removeItem('user')
        } else {
          sessionStorage.setItem('token', action.payload.token)
          sessionStorage.setItem('user', JSON.stringify(action.payload.user))
          localStorage.removeItem('token')
          localStorage.removeItem('user')
        }
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string
      })
      .addCase(fetchMe.fulfilled, (state, action) => {
        state.user = action.payload
        localStorage.setItem('user', JSON.stringify(action.payload))
      })
  },
})

export const { logout, setUser } = authSlice.actions
export default authSlice.reducer
