import { Box, Button, Drawer, List, ListItem, ListItemButton, ListItemText, TextField, Typography } from "@mui/material"
import { forwardRef, useEffect, useState } from "react"
import { getHosts, sayHello, deleteHost } from "../util"
import { Host } from "../model"

export const Hosts = forwardRef<HTMLDivElement>((props, ref): JSX.Element => {

    const token = localStorage.getItem("JWT")
    const [hosts, setHosts] = useState<Host[]>([])
    const [remoteFqdn, setRemoteFqdn] = useState('')

    const [selectedHost, setSelectedHost] = useState<Host | null>(null)
    const [newRole, setNewRole] = useState<string>('')
    const [newScore, setNewScore] = useState<number>(0)

    useEffect(() => {
        getHosts().then(setHosts)
    }, [])

    return (
        <div ref={ref} {...props}>
            <Box sx={{ position: 'absolute', width: '100%' }}>
                <Box sx={{ display: 'flex', gap: '10px' }}>
                    <TextField
                        label="remote fqdn"
                        variant="outlined"
                        value={remoteFqdn}
                        sx={{ flexGrow: 1 }}
                        onChange={(e) => {
                            setRemoteFqdn(e.target.value)
                        }}
                    />
                    <Button
                        variant="contained"
                        onClick={(_) => {
                            if (!token) return
                            sayHello(token, remoteFqdn)
                        }}
                    >
                        GO
                    </Button>
                </Box>
                <Typography>Hosts</Typography>
                <List
                    disablePadding
                >
                    {hosts.map((host) => (
                        <ListItem key={host.ccaddr}
                            disablePadding
                        >
                            <ListItemButton
                                onClick={() => {
                                    setNewRole(host.role)
                                    setNewScore(host.score)
                                    setSelectedHost(host)
                                }}
                            >
                                <ListItemText primary={host.fqdn} secondary={`${host.ccaddr}`} />
                                <ListItemText>{`${host.role}(${host.score})`}</ListItemText>
                            </ListItemButton>
                        </ListItem>
                    ))}
                </List>
            </Box>
            <Drawer
                anchor="right"
                open={selectedHost !== null}
                onClose={() => {
                    setSelectedHost(null)
                }}
            >
                <Box
                    width="50vw"
                    display="flex"
                    flexDirection="column"
                    gap={1}
                    padding={2}
                >
                    <Typography>{selectedHost?.ccaddr}</Typography>
                    <pre>{JSON.stringify(selectedHost, null, 2)}</pre>
                    <TextField
                        label="new role"
                        variant="outlined"
                        value={newRole}
                        sx={{ flexGrow: 1 }}
                        onChange={(e) => {
                            setNewRole(e.target.value)
                        }}
                    />
                    <TextField
                        label="new score"
                        variant="outlined"
                        value={newScore}
                        sx={{ flexGrow: 1 }}
                        onChange={(e) => {
                            setNewScore(Number(e.target.value))
                        }}
                    />
                    <Button
                        variant="contained"
                        onClick={(_) => {
                            if (!token) return
                        }}
                    >
                        Update
                    </Button>
                    <Button
                        variant="contained"
                        onClick={(_) => {
                            if (!token) return
                            if (!selectedHost) return
                            deleteHost(token, selectedHost.fqdn)
                            setSelectedHost(null)
                        }}
                        color="error"
                    >
                        Delete
                    </Button>
                </Box>
            </Drawer>
        </div>
    )
})

Hosts.displayName = "Hosts"

