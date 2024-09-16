package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Usage: kube-pfw <namespace>")
	}

	namespace := os.Args[1]

	if err := checkKubectlExists(); err != nil {
		return fmt.Errorf("kubectl check failed: %w", err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	if len(services.Items) == 0 {
		return fmt.Errorf("no services found in namespace %s", namespace)
	}

	selectedService, err := selectService(services.Items)
	if err != nil {
		return fmt.Errorf("service selection failed: %w", err)
	}

	selectedPort, err := selectPort(selectedService)
	if err != nil {
		return fmt.Errorf("port selection failed: %w", err)
	}

	localPort, err := getLocalPort()
	if err != nil {
		return fmt.Errorf("local port input failed: %w", err)
	}

	if err := runPortForward(selectedService, selectedPort, localPort, namespace); err != nil {
		return fmt.Errorf("port-forward failed: %w", err)
	}

	return nil
}

func checkKubectlExists() error {
	_, err := exec.LookPath("kubectl")
	if err != nil {
		return fmt.Errorf("kubectl is not installed or not in PATH")
	}
	return nil
}

func selectService(services []corev1.Service) (*corev1.Service, error) {
	fmt.Println("* service:")
	for i, svc := range services {
		ports := []string{}
		for _, port := range svc.Spec.Ports {
			ports = append(ports, fmt.Sprintf("%d", port.Port))
		}
		fmt.Printf("  %d. %s ( port %s )\n", i+1, svc.Name, strings.Join(ports, " , "))
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the number: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index < 1 || index > len(services) {
		return nil, fmt.Errorf("invalid selection")
	}

	return &services[index-1], nil
}

func selectPort(service *corev1.Service) (int32, error) {
	if len(service.Spec.Ports) == 1 {
		return service.Spec.Ports[0].Port, nil
	}

	fmt.Printf("* %s:\n", service.Name)
	for i, port := range service.Spec.Ports {
		fmt.Printf("  %d. %d\n", i+1, port.Port)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the number: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("failed to read input: %w", err)
	}

	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index < 1 || index > len(service.Spec.Ports) {
		return 0, fmt.Errorf("invalid selection")
	}

	return service.Spec.Ports[index-1].Port, nil
}

func getLocalPort() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("* Local Port: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	localPort := strings.TrimSpace(input)
	if _, err := strconv.Atoi(localPort); err != nil {
		return "", fmt.Errorf("invalid port number")
	}

	return localPort, nil
}

func runPortForward(service *corev1.Service, remotePort int32, localPort, namespace string) error {
	cmd := exec.Command("kubectl", "port-forward", fmt.Sprintf("service/%s", service.Name), fmt.Sprintf("%s:%d", localPort, remotePort), "-n", namespace)

	fmt.Printf("Exec Command: %s\n", cmd.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
