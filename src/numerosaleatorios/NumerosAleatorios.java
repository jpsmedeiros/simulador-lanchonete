/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package numerosaleatorios;

import java.util.Random;

/**
 *
 * @author João Pedro
 */
public class NumerosAleatorios {

    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) {
        // TODO code application logic here
        double A = 234563439d;
        double M = Math.pow(2, 48);
        double C = 457842;
        double x0 = 3;
        double qtdNumeros = 10000;
        double ultimoValor = x0;
        double xAtual, resultado;
        //System.out.println("========== NUMEROS GERADOS PELO METODO ==========");
        for(int i=0; i<qtdNumeros; i++){
            xAtual = ((A * ultimoValor) % M);
            resultado = xAtual / M;
            ultimoValor =  xAtual;
            //System.out.printf("%.1f\n", resultado);
        }
        //System.out.println("========== NUMEROS GERADOS PELO SISTEMA ==========");
        Random r = new Random();
        double randomValue;
        for(int i=0; i<qtdNumeros; i++){
            randomValue = 0 + (1 - 0) * r.nextDouble();
            System.out.printf("%.1f\n", randomValue);
        }
    }
    
}
