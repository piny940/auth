import type { Metadata } from 'next'
import localFont from 'next/font/local'
import './globals.css'
import { AppBar, Box, Container, Toolbar, Typography } from '@mui/material'
import Link from 'next/link'

const geistSans = localFont({
  src: './fonts/GeistVF.woff',
  variable: '--font-geist-sans',
  weight: '100 900',
})
const geistMono = localFont({
  src: './fonts/GeistMonoVF.woff',
  variable: '--font-geist-mono',
  weight: '100 900',
})

export const metadata: Metadata = {
  title: 'mikan login',
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="ja">
      <body className={`${geistSans.variable} ${geistMono.variable}`}>
        <Box>
          <AppBar position="sticky">
            <Toolbar>
              <Link href="/" className="unstyled">
                <Typography variant="h5" fontWeight="bold">
                  mikan
                </Typography>
              </Link>
            </Toolbar>
          </AppBar>
          {children}
        </Box>
      </body>
    </html>
  )
}
