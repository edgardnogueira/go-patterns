package proxy

import (
	"fmt"
	"strings"
)

// User represents a system user with permissions.
type User struct {
	Username string
	Role     string
}

// ProtectionProxy controls access to the image based on user permissions.
// It adds a security layer to allow or deny operations based on authentication and authorization.
type ProtectionProxy struct {
	BaseProxy
	user              *User
	allowedRoles      map[string]bool
	allowedExtensions map[string]bool
}

// NewProtectionProxy creates a new protection proxy.
func NewProtectionProxy(realImage Image, user *User) *ProtectionProxy {
	proxy := &ProtectionProxy{
		user: user,
		allowedRoles: map[string]bool{
			"admin":  true,
			"editor": true,
		},
		allowedExtensions: map[string]bool{
			".jpg":  true,
			".png":  true,
			".gif":  true,
			".bmp":  true,
			".webp": true,
		},
	}
	proxy.realImage = realImage
	return proxy
}

// Display checks user permissions before allowing display of the image.
func (p *ProtectionProxy) Display() error {
	if err := p.checkAccess("display"); err != nil {
		return err
	}
	
	if err := p.validateFile(); err != nil {
		return err
	}
	
	fmt.Println("Protection proxy allowing display operation")
	return p.realImage.Display()
}

// GetFilename returns the image's filename if the user has permission.
func (p *ProtectionProxy) GetFilename() string {
	if err := p.checkAccess("getFilename"); err != nil {
		fmt.Printf("Access denied: %v\n", err)
		return "Access Denied"
	}
	return p.realImage.GetFilename()
}

// GetWidth returns the width if the user has permission.
func (p *ProtectionProxy) GetWidth() int {
	if err := p.checkAccess("getMetadata"); err != nil {
		fmt.Printf("Access denied: %v\n", err)
		return 0
	}
	return p.realImage.GetWidth()
}

// GetHeight returns the height if the user has permission.
func (p *ProtectionProxy) GetHeight() int {
	if err := p.checkAccess("getMetadata"); err != nil {
		fmt.Printf("Access denied: %v\n", err)
		return 0
	}
	return p.realImage.GetHeight()
}

// GetSize returns the size if the user has permission.
func (p *ProtectionProxy) GetSize() int64 {
	if err := p.checkAccess("getMetadata"); err != nil {
		fmt.Printf("Access denied: %v\n", err)
		return 0
	}
	return p.realImage.GetSize()
}

// GetMetadata returns the metadata if the user has permission.
func (p *ProtectionProxy) GetMetadata() map[string]string {
	if err := p.checkAccess("getMetadata"); err != nil {
		fmt.Printf("Access denied: %v\n", err)
		return map[string]string{"error": "Access Denied"}
	}
	return p.realImage.GetMetadata()
}

// checkAccess checks if the user has permission for the operation.
func (p *ProtectionProxy) checkAccess(operation string) error {
	if p.user == nil {
		return fmt.Errorf("authentication required: no user provided")
	}
	
	// Check if user has an allowed role
	if !p.allowedRoles[p.user.Role] {
		return fmt.Errorf("authorization failed: role '%s' does not have permission for operation '%s'", 
			p.user.Role, operation)
	}
	
	fmt.Printf("Access granted to %s (%s) for operation '%s'\n", p.user.Username, p.user.Role, operation)
	return nil
}

// validateFile checks if the file has an allowed extension.
func (p *ProtectionProxy) validateFile() error {
	filename := p.realImage.GetFilename()
	
	// Check file extension
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex == -1 {
		return fmt.Errorf("invalid file: missing extension")
	}
	
	extension := strings.ToLower(filename[dotIndex:])
	if !p.allowedExtensions[extension] {
		return fmt.Errorf("unsupported file type: %s", extension)
	}
	
	return nil
}

// SetAllowedRoles updates the list of roles that are allowed to access the image.
func (p *ProtectionProxy) SetAllowedRoles(roles []string) {
	p.allowedRoles = make(map[string]bool)
	for _, role := range roles {
		p.allowedRoles[role] = true
	}
}

// SetAllowedExtensions updates the list of file extensions that are allowed.
func (p *ProtectionProxy) SetAllowedExtensions(extensions []string) {
	p.allowedExtensions = make(map[string]bool)
	for _, ext := range extensions {
		// Ensure extensions start with a dot
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		p.allowedExtensions[strings.ToLower(ext)] = true
	}
}
